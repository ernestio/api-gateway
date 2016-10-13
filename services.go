/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

// getServicesHandler : responds to GET /services/ with a list of all
// services for current user group
func getServicesHandler(c echo.Context) (err error) {
	var services []Service
	var list []Service
	var body []byte
	var service Service

	au := authenticatedUser(c)
	service.FindAll(au, &services)
	for _, s := range services {
		exists := false
		for i, e := range list {
			if e.Name == s.Name {
				if e.Version.Before(s.Version) {
					list[i] = s
				}
				exists = true
			}
		}
		if exists == false {
			list = append(list, s)
		}
	}

	if body, err = json.Marshal(list); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// getServiceBuildsHandler : gets the list of builds for the specified
// service
func getServiceBuildsHandler(c echo.Context) error {
	au := authenticatedUser(c)

	query := getParamFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	list, err := getServicesOutput(query)
	if err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	return c.JSON(http.StatusOK, list)
}

// getServiceHandler : responds to GET /services/:service with the
// details of an existing service
func getServiceHandler(c echo.Context) (err error) {
	var s Service
	var services []Service
	var o ServiceRender

	au := authenticatedUser(c)
	query := getParamFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if err = s.Find(query, &services); err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	if len(services) > 0 {
		o.Render(services[0])
		if body, err := o.ToJson(); err != nil {
			return c.JSONBlob(500, []byte(err.Error()))
		} else {
			return c.JSONBlob(http.StatusOK, body)
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

// getServiceBuildHandler : gets the details of a specific service build
func getServiceBuildHandler(c echo.Context) (err error) {
	var list []ServiceRender

	au := authenticatedUser(c)
	query := getParamFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if list, err = getServicesOutput(query); err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	if len(list) > 0 {
		return c.JSON(http.StatusOK, list[0])
	}
	return c.JSON(http.StatusNotFound, nil)
}

// TODO : WTF is this doing??
func searchServicesHandler(c echo.Context) error {
	au := authenticatedUser(c)

	query := getSearchFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	list, err := getServicesOutput(query)
	if err != nil {
		return ErrInternal
	}

	return c.JSON(http.StatusOK, list)
}

// resetServiceHandler : Respons to POST /services/:service/reset/ and updates the
// service status to errored from in_progress
func resetServiceHandler(c echo.Context) error {
	var s Service
	var services []Service

	name := c.Param("service")

	au := authenticatedUser(c)
	filter := make(map[string]interface{})
	filter["group_id"] = au.GroupID
	filter["name"] = name
	if err := s.Find(filter, &services); err != nil {
		log.Println(err.Error())
		return c.JSONBlob(500, []byte("Internal Error"))
	}

	if len(services) == 0 {
		return c.JSONBlob(404, []byte("Service not found with this name"))
	}

	s = services[0]

	if s.Status != "in_progress" {
		return c.JSONBlob(200, []byte("Reset only applies to 'in progress' serices, however service '"+name+"' is on status '"+s.Status))
	}

	if err := s.Reset(); err != nil {
		log.Println(err.Error())
		return c.JSONBlob(500, []byte("Internal error"))
	}

	return c.String(200, "success")
}

func createUUIDHandler(c echo.Context) error {
	var s struct {
		ID string `json:"id"`
	}
	req := c.Request()
	body, err := ioutil.ReadAll(req.Body())
	if err != nil {
		return c.JSONBlob(500, []byte("Invalid input"))
	}

	json.Unmarshal(body, &s)
	id := generateStreamID(s.ID)

	return c.JSONBlob(http.StatusOK, []byte(`{"uuid":"`+id+`"}`))
}

// createServiceHandler : Will receive a service application
func createServiceHandler(c echo.Context) error {
	var s ServiceInput
	var err error
	var body []byte
	var definition []byte
	var datacenter []byte
	var group []byte

	payload := ServicePayload{}
	au := authenticatedUser(c)

	if au.GroupID == 0 {
		return c.JSONBlob(401, []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action"))
	}
	if s, definition, body, err = mapInputService(c); err != nil {
		return c.JSONBlob(400, []byte(err.Error()))
	}
	payload.Service = (*json.RawMessage)(&body)

	// Get datacenter
	if datacenter, err = getDatacenter(s.Datacenter, au.GroupID); err != nil {
		return c.JSONBlob(404, []byte(err.Error()))
	}
	payload.Datacenter = (*json.RawMessage)(&datacenter)

	// Get group
	if group, err = getGroup(au.GroupID); err != nil {
		return c.JSONBlob(http.StatusNotFound, []byte(err.Error()))
	}
	payload.Group = (*json.RawMessage)(&group)

	// Generate service ID
	payload.ID = generateServiceID(s.Name + "-" + s.Datacenter)

	// Get previous service if exists
	if previous, err := getService(s.Name, au.GroupID); err != nil {
		return c.JSONBlob(http.StatusNotFound, []byte(err.Error()))
	} else {
		if previous != nil {
			payload.PrevID = previous.ID
			if previous.Status == "in_progress" {
				return c.JSONBlob(http.StatusNotFound, []byte(`"Your service process is 'in progress' if your're sure you want to fix it please reset it first"`))
			}
		}
	}

	var service []byte
	if service, err = mapCreateDefinition(payload); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	var datacenterStruct struct {
		ID   uint   `json:"id"`
		Type string `json:"type"`
	}
	json.Unmarshal(datacenter, &datacenterStruct)

	datacenterID, _ := strconv.Atoi(string(datacenterStruct.ID))

	ss := Service{
		ID:           payload.ID,
		Name:         s.Name,
		Type:         datacenterStruct.Type,
		GroupID:      au.GroupID,
		DatacenterID: datacenterID,
		Version:      time.Now(),
		Status:       "in_progress",
		Definition:   string(definition),
		Maped:        string(service),
	}
	if err := ss.Save(); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	// Apply changes
	n.Publish("service.create", service)

	return c.JSONBlob(http.StatusOK, []byte(`{"id":"`+payload.ID+`"}`))
}

func updateServiceHandler(c echo.Context) error {
	return echo.NewHTTPError(405, "Not implemented")
}

// Deletes a service by name
func deleteServiceHandler(c echo.Context) error {
	var raw []byte
	var err error

	au := authenticatedUser(c)

	if raw, err = getServiceRaw(c.Param("name"), au.GroupID); err != nil {
		return echo.NewHTTPError(404, err.Error())
	}

	s := Service{}
	json.Unmarshal(raw, &s)

	if s.Status == "in_progress" {
		return c.JSONBlob(400, []byte(`"Service is already applying some changes, please wait until they are done"`))
	}

	query := []byte(`{"previous_id":"` + s.ID + `","datacenter":{"type":"` + s.Type + `"}}`)
	if msg, err := n.Request("definition.map.deletion", query, 1*time.Second); err != nil {
		return c.JSONBlob(500, []byte(`"Couldn't map the service"`))
	} else {
		n.Publish("service.delete", msg.Data)
	}

	parts := strings.Split(s.ID, "-")
	stream := parts[len(parts)-1]

	return c.JSONBlob(http.StatusOK, []byte(`{"id":"`+s.ID+`","stream_id":"`+stream+`"}`))
}
