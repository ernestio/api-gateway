/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Service holds the service response from service-store
type Service struct {
	ID           int                    `json:"id"`
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
		query = fmt.Sprintf(`{"name": "%s"}`, c.Param("service"))
	} else {
		query = fmt.Sprintf(`{"name": "%s", "group_id": %d}`, c.Param("service"), au.GroupID)
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

func createServiceHandler(c echo.Context) error {
	return ErrNotImplemented
}

func updateServiceHandler(c echo.Context) error {
	return ErrNotImplemented
}

func deleteServiceHandler(c echo.Context) error {
	return ErrNotImplemented
}
