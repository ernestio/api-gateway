/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

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
	if au.GroupID == 0 {
		return c.JSONBlob(401, []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action"))
	}
	d.GroupID = au.GroupID

	data, err := json.Marshal(d)
	if err != nil {
		return ErrInternal
	}

	// Does the datacenter already exist
	msg, err := n.Request("datacenter.get", []byte(`{"name":"`+d.Name+`"}`), 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}
	if string(msg.Data) != `{"error":"not found"}` {
		return c.JSONBlob(409, []byte("Datacenter name already in use"))
	}

	msg, err = n.Request("datacenter.set", data, 1*time.Second)
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
