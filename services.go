/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
)

// Service holds the service response from service-store
type Service struct {
	ID           string                 `json:"id"`
	GroupID      int                    `json:"group_id"`
	DatacenterID int                    `json:"datacenter_id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Version      string                 `json:"version"`
	Options      map[string]interface{} `json:"options"`
	Status       string                 `json:"status"`
	Endpoint     string                 `json:"endpoint"`
	Definition   interface{}            `json:"definition"`
}

// Validate the service
func (d *Service) Validate() error {
	if d.Name == "" {
		return errors.New("Service name is empty")
	}

	if d.DatacenterID == 0 {
		return errors.New("Service group is empty")
	}

	if d.Type == "" {
		return errors.New("Service type is empty")
	}

	return nil
}

// Map : maps a service from a request's body and validates the input
func (d *Service) Map(c echo.Context) *echo.HTTPError {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return ErrBadReqBody
	}

	err = d.Validate()
	if err != nil {
		return ErrBadReqBody
	}

	return nil
}

func getServicesHandler(c echo.Context) (err error) {
	var list []OutputService
	au := authenticatedUser(c)

	query := getSearchFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if list, err = getServicesOutput(query); err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	return c.JSON(http.StatusOK, list)
}

func getServiceBuildsHandler(c echo.Context) error {
	// get the service name
	au := authenticatedUser(c)

	// Get all builds for service name
	qb := getParamFilter(c)
	if au.Admin != true {
		qb["group_id"] = au.GroupID
	}

	list, err := getServicesOutput(qb)
	if err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	return c.JSON(http.StatusOK, list)
}

func getServiceHandler(c echo.Context) (err error) {
	var list []OutputService

	au := authenticatedUser(c)
	query := getParamFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if list, err = getServicesOutput(query); err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	return c.JSON(http.StatusOK, list[0])
}

func getServiceBuildHandler(c echo.Context) (err error) {
	var list []OutputService

	au := authenticatedUser(c)
	query := getParamFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if list, err = getServicesOutput(query); err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	return c.JSON(http.StatusOK, list[0])
}

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

func resetServiceHandler(c echo.Context) error {
	au := authenticatedUser(c)
	if status, err := resetService(au, c.Param("service")); err != nil {
		return c.JSONBlob(status, []byte(err.Error()))
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
	var datacenter []byte
	var group []byte
	var action = "service.create"

	payload := ServicePayload{}
	au := authenticatedUser(c)

	if s, body, err = mapInputService(c); err != nil {
		return c.JSONBlob(400, []byte(err.Error()))
	}
	payload.Service = (*json.RawMessage)(&body)

	// Get datacenter
	if datacenter, err = getDatacenter(c.Param("datacenter"), au.GroupID, s.Provider); err != nil {
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
			if previous.Status == "errored" {
				action = "service.patch"
			}
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
		Type string `json:"typ"`
	}
	json.Unmarshal(datacenter, &datacenterStruct)

	version := "0.1.1"
	status := "in_progress"
	options := "{}"
	definition := string(body)
	mapping := string(service)
	saveService(payload.ID, s.Name, datacenterStruct.Type, version, status, options, definition, mapping, uint(au.GroupID), datacenterStruct.ID)

	// Apply changes
	n.Publish(action, service)

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

	query := []byte(`{"previous_id":"` + s.ID + `"}`)
	if msg, err := n.Request("definition.map.deletion", query, 1*time.Second); err != nil {
		return c.JSONBlob(500, []byte(`"Couldn't map the service"`))
	} else {
		n.Publish("service.delete", msg.Data)
	}

	parts := strings.Split(s.ID, "-")
	stream := parts[len(parts)-1]

	return c.JSONBlob(http.StatusOK, []byte(`{"id":"`+s.ID+`","stream_id":"`+stream+`"}`))
}
