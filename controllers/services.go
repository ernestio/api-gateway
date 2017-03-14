/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
	"github.com/labstack/echo"
)

// GetServicesHandler : responds to GET /services/ with a list of all
// services for current user group
func GetServicesHandler(c echo.Context) (err error) {
	var services []models.Service
	var list []models.Service
	var body []byte
	var service models.Service
	var user models.User

	users := user.FindAllKeyValue()

	au := AuthenticatedUser(c)
	if err := service.FindAll(au, &services); err != nil {
		log.Println(err)
	}
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
			for id, name := range users {
				if id == s.UserID {
					s.UserName = name
				}
			}
			list = append(list, s)
		}
	}

	if body, err = json.Marshal(list); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// GetServiceBuildsHandler : gets the list of builds for the specified
// service
func GetServiceBuildsHandler(c echo.Context) error {
	var user models.User

	users := user.FindAllKeyValue()
	au := AuthenticatedUser(c)

	query := h.GetParamFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	list, err := getServicesOutput(query)
	if err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}
	for i := range list {
		for id, name := range users {
			if id == list[i].UserID {
				list[i].UserName = name
			}
		}
	}

	return c.JSON(http.StatusOK, list)
}

// GetServiceHandler : responds to GET /services/:service with the
// details of an existing service
func GetServiceHandler(c echo.Context) (err error) {
	var s models.Service
	var services []models.Service
	var o views.ServiceRender
	var body []byte

	au := AuthenticatedUser(c)
	query := h.GetParamFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if err = s.Find(query, &services); err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	if len(services) > 0 {
		if err := o.Render(services[0]); err != nil {
			log.Println(err)
			return err
		}
		if body, err = o.ToJSON(); err != nil {
			return c.JSONBlob(500, []byte(err.Error()))
		}
		return c.JSONBlob(http.StatusOK, body)
	}
	return c.JSON(http.StatusNotFound, nil)
}

// GetServiceBuildHandler : gets the details of a specific service build
func GetServiceBuildHandler(c echo.Context) (err error) {
	var list []views.ServiceRender

	au := AuthenticatedUser(c)
	query := h.GetParamFilter(c)
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

// SearchServicesHandler : Finds all services
func SearchServicesHandler(c echo.Context) error {
	au := AuthenticatedUser(c)

	query := h.GetSearchFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	list, err := getServicesOutput(query)
	if err != nil {
		return h.ErrInternal
	}

	return c.JSON(http.StatusOK, list)
}

// ResetServiceHandler : Respons to POST /services/:service/reset/ and updates the
// service status to errored from in_progress
func ResetServiceHandler(c echo.Context) error {
	var s models.Service
	var services []models.Service

	name := c.Param("service")

	au := AuthenticatedUser(c)
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

// CreateUUIDHandler : Creates an unique id
func CreateUUIDHandler(c echo.Context) error {
	var s struct {
		ID string `json:"id"`
	}
	req := c.Request()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return c.JSONBlob(500, []byte("Invalid input"))
	}

	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return err
	}
	id := generateStreamID(s.ID)

	return c.JSONBlob(http.StatusOK, []byte(`{"uuid":"`+id+`"}`))
}

// CreateServiceHandler : Will receive a service application
func CreateServiceHandler(c echo.Context) error {
	var s ServiceInput
	var err error
	var body []byte
	var definition []byte
	var datacenter []byte
	var group []byte
	var previous *models.Service

	payload := ServicePayload{}
	au := AuthenticatedUser(c)

	if au.GroupID == 0 {
		body := "Current user does not belong to any group."
		body += "\nPlease assign the user to a group before performing this action"
		return c.JSONBlob(401, []byte(body))
	}

	// Parse the input service as usual
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
	var currentUser models.User
	if err := currentUser.FindByUserName(au.Username, &currentUser); err != nil {
		log.Println(err)
		return err
	}

	// Generate service ID
	payload.ID = generateServiceID(s.Name + "-" + s.Datacenter)

	// Get previous service if exists
	if previous, err = getService(s.Name, au.GroupID); err != nil {
		return c.JSONBlob(http.StatusNotFound, []byte(err.Error()))
	}

	if previous != nil {
		payload.PrevID = previous.ID
		if previous.Status == "in_progress" {
			return c.JSONBlob(http.StatusNotFound, []byte(`"Your service process is 'in progress' if your're sure you want to fix it please reset it first"`))
		}
	}

	var service []byte
	isAnImport := strings.Contains(c.Path(), "/import/")

	if body, err = json.Marshal(payload); err != nil {
		return h.ErrInternal
	}
	var def models.Definition
	if isAnImport == true {
		service, err = def.MapImport(body)
	} else {
		service, err = def.MapCreation(body)
	}

	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	var datacenterStruct struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(datacenter, &datacenterStruct); err != nil {
		log.Println(err)
		return err
	}

	ss := models.Service{
		ID:           payload.ID,
		Name:         s.Name,
		Type:         datacenterStruct.Type,
		GroupID:      au.GroupID,
		UserID:       currentUser.ID,
		DatacenterID: datacenterStruct.ID,
		Version:      time.Now(),
		Status:       "in_progress",
		Definition:   string(definition),
		Maped:        string(service),
	}

	if err := ss.Save(); err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	// Apply changes
	if isAnImport == true {
		err = ss.RequestImport(service)
	} else {
		err = ss.RequestCreation(service)
	}

	if err != nil {
		h.L.Error(err.Error())
		return err
	}

	return c.JSONBlob(http.StatusOK, []byte(`{"id":"`+payload.ID+`"}`))
}

// UpdateServiceHandler : Not implemented
func UpdateServiceHandler(c echo.Context) error {
	h.L.Warning("UpdateServiceHandler not implemented")
	return echo.NewHTTPError(405, "Not implemented")
}

// DeleteServiceHandler : Deletes a service by name
func DeleteServiceHandler(c echo.Context) error {
	var raw []byte
	var err error
	var def models.Definition

	au := AuthenticatedUser(c)

	if raw, err = getServiceRaw(c.Param("name"), au.GroupID); err != nil {
		return echo.NewHTTPError(404, err.Error())
	}

	s := models.Service{}
	if err := json.Unmarshal(raw, &s); err != nil {
		log.Println(err)
		return err
	}

	if s.Status == "in_progress" {
		return c.JSONBlob(400, []byte(`"Service is already applying some changes, please wait until they are done"`))
	}

	body, err := def.MapDeletion(s.ID, s.Type)
	if err != nil {
		return c.JSONBlob(500, []byte(`"Couldn't map the service"`))
	}
	if err := s.RequestDeletion(body); err != nil {
		return c.JSONBlob(500, []byte(`"Couldn't call service.delete"`))
	}

	parts := strings.Split(s.ID, "-")
	stream := parts[len(parts)-1]

	return c.JSONBlob(http.StatusOK, []byte(`{"id":"`+s.ID+`","stream_id":"`+stream+`"}`))
}

// ForceServiceDeletionHandler : Deletes a service by name forcing it
func ForceServiceDeletionHandler(c echo.Context) error {
	var raw []byte
	var err error
	var service models.Service

	au := AuthenticatedUser(c)

	if raw, err = getServiceRaw(c.Param("name"), au.GroupID); err != nil {
		return echo.NewHTTPError(404, err.Error())
	}

	s := models.Service{}
	if err := json.Unmarshal(raw, &s); err != nil {
		log.Println(err)
		return echo.NewHTTPError(500, err.Error())
	}

	if err := service.DeleteByName(c.Param("name")); err != nil {
		log.Println(err)
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSONBlob(http.StatusOK, []byte(`{"id":"`+s.ID+`"}`))
}
