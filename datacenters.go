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

// Datacenter holds the datacenter response from datacenter-store
type Datacenter struct {
	ID              int    `json:"id"`
	GroupID         int    `json:"group_id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Region          string `json:"region"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	VCloudURL       string `json:"vcloud_url"`
	VseURL          string `json:"vse_url"`
	ExternalNetwork string `json:"external_network"`
}

// Validate the datacenter
func (d *Datacenter) Validate() error {
	if d.Name == "" {
		return errors.New("Datacenter name is empty")
	}

	if d.GroupID == 0 {
		return errors.New("Datacenter group is empty")
	}

	if d.Type == "" {
		return errors.New("Datacenter type is empty")
	}

	if d.Username == "" {
		return errors.New("Datacenter username is empty")
	}

	if d.Password == "" {
		return errors.New("Datacenter password is empty")
	}

	if d.Type == "vcloud" && d.VCloudURL == "" {
		return errors.New("Datacenter vcloud url is empty")
	}

	return nil
}

// Map : maps a datacenter from a request's body and validates the input
func (d *Datacenter) Map(c echo.Context) *echo.HTTPError {
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

// getDatacentersHandler : get all datacenters
func getDatacentersHandler(c echo.Context) error {
	au := authenticatedUser(c)
	query := fmt.Sprintf(`{"group_id": %d}`, au.GroupID)
	msg, err := n.Request("datacenter.find", []byte(query), 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

// getDatacenterHandler : get a datancenter by id
func getDatacenterHandler(c echo.Context) error {
	var query string
	au := authenticatedUser(c)

	if au.Admin {
		query = fmt.Sprintf(`{"id": %s}`, c.Param("datacenter"))
	} else {
		query = fmt.Sprintf(`{"id": %s, "group_id": %d}`, c.Param("datacenter"), au.GroupID)
	}

	msg, err := n.Request("datacenter.get", []byte(query), 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

// createDatacenterHandler : Endpoint to create a datacenter
func createDatacenterHandler(c echo.Context) error {
	var d Datacenter
	if d.Map(c) != nil {
		return ErrBadReqBody
	}

	au := authenticatedUser(c)
	if au.Admin != true && d.GroupID != au.GroupID {
		return ErrUnauthorized
	}

	data, err := json.Marshal(d)
	if err != nil {
		return ErrInternal
	}

	msg, err := n.Request("datacenter.set", data, 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

// updateDatacenterHandler : Updates a datacenter through its store
func updateDatacenterHandler(c echo.Context) error {
	var d Datacenter
	if d.Map(c) != nil {
		return ErrBadReqBody
	}

	au := authenticatedUser(c)
	if au.Admin != true || d.GroupID != au.GroupID {
		return ErrUnauthorized
	}

	data, err := json.Marshal(d)
	if err != nil {
		return ErrInternal
	}

	msg, err := n.Request("datacenter.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

// deleteDatacenterHandler : Deletes a datancenter though its store
func deleteDatacenterHandler(c echo.Context) error {
	au := authenticatedUser(c)

	query := fmt.Sprintf(`{"id": %s, "group_id": %d}`, c.Param("datacenter"), au.GroupID)
	msg, err := n.Request("datacenter.del", []byte(query), 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.String(http.StatusOK, "")
}
