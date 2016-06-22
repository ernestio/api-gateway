/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/labstack/echo"
	"github.com/nats-io/nats"
	"github.com/nu7hatch/gouuid"
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

	if d.GroupID == 0 {
		return errors.New("Service group is empty")
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

func getServicesHandler(c echo.Context) error {
	msg, err := n.Request("service.find", nil, 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getServiceHandler(c echo.Context) error {
	var query string
	au := authenticatedUser(c)

	if au.Admin {
		query = fmt.Sprintf(`{"id": "%s"}`, c.Param("service"))
	} else {
		query = fmt.Sprintf(`{"id": "%s", "group_id": %d}`, c.Param("service"), au.GroupID)
	}

	msg, err := n.Request("service.get", []byte(query), 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

// createServiceHandler : Will receive a service application
func createServiceHandler(c echo.Context) error {
	var msg *nats.Msg

	type Service struct {
		Datacenter string `json:"datacenter"`
		Provider   string `json:"provider"`
		Name       string `json:"name"`
	}
	au := authenticatedUser(c)
	req := c.Request()
	body, err := ioutil.ReadAll(req.Body())

	// Normalize input body to json
	ctype := http.DetectContentType(body)
	if ctype != "application/json" && ctype != "application/yaml" {
		return echo.NewHTTPError(400, "Invalid input format")
	} else if ctype == "application/yaml" {
		body, err = yaml.JSONToYAML(body)
	}

	// TODO check content type
	s := Service{}
	if err = json.Unmarshal(body, &s); err != nil {
		return echo.NewHTTPError(400, "Invalid input")
	}

	// Get datacenter
	query := fmt.Sprintf(`{"id": %s, "group_id": %d}`, c.Param("datacenter"), au.GroupID)
	if msg, err = n.Request("datacenter.find", []byte(query), 1*time.Second); err != nil {
		return ErrGatewayTimeout
	}
	if strings.Contains(string(msg.Data), `"error"`) {
		return echo.NewHTTPError(http.StatusNotFound, "Specified datacenter does not exist")
	}

	// TODO override datacenter if is a fake one
	if s.Provider == "fake" {
	}
	datacenter := msg.Data

	// FIXME : Is still needed this block about groups?
	// get client data
	query = fmt.Sprintf(`{"id": %d}`, au.GroupID)
	msg, err = n.Request("group.get", []byte(query), 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}
	if strings.Contains(string(msg.Data), `"error"`) {
		return echo.NewHTTPError(http.StatusNotFound, "Specified group does not exist")
	}
	group := msg.Data

	// Calculate the id
	sufix := md5.Sum([]byte(s.Name + "-" + s.Datacenter))
	prefix, err := uuid.NewV4()
	serviceID := prefix.String() + "-" + string(sufix[:])

	// We need the status of previous service
	query = fmt.Sprintf(`{"name":"%s","group_id":"%d"}`, s.Name, au.GroupID)
	if msg, err = n.Request("service.get", []byte(query), 1*time.Second); err != nil {
		return ErrGatewayTimeout
	}
	var p struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	subject := "service.create"
	if !strings.Contains(string(msg.Data), `{"error":"not found"}`) {
		json.Unmarshal(msg.Data, &p)
		if p.Status == "errored" {
			subject = "service.patch"
		}
	}

	// MAP definition
	var payload struct {
		ID         string      `json:"id"`
		PrevID     string      `json:"previous_id"`
		Datacenter interface{} `json:"datacenter"`
		Group      interface{} `json:"client"`
		Service    interface{} `json:"service"`
	}

	payload.ID = serviceID
	payload.Service = body
	payload.Datacenter = datacenter
	payload.Group = group
	payload.PrevID = p.ID

	var payloadBody []byte
	if payloadBody, err = json.Marshal(payload); err != nil {
		return echo.NewHTTPError(400, "Provided yaml is not valid")
	}
	if msg, err = n.Request("definition.map_create", []byte(payloadBody), 1*time.Second); err != nil {
		return echo.NewHTTPError(400, "Provided yaml is not valid")
	}
	mappedService := msg.Data

	// Apply changes
	n.Publish(subject, mappedService)

	// Respond with { id: id } -------->
	return c.JSONBlob(http.StatusOK, []byte(`{"id":"`+serviceID+`"}`))
}

func updateServiceHandler(c echo.Context) error {
	return ErrNotImplemented
}

func deleteServiceHandler(c echo.Context) error {
	return ErrNotImplemented
}
