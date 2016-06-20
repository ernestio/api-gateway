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

// Group holds the group response from group-store
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Validate the group
func (g *Group) Validate() error {
	if g.Name == "" {
		return errors.New("Group name is empty")
	}

	return nil
}

// Map : maps a group from a request's body and validates the input
func (g *Group) Map(c echo.Context) *echo.HTTPError {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &g)
	if err != nil {
		return ErrBadReqBody
	}

	err = g.Validate()
	if err != nil {
		return ErrBadReqBody
	}

	return nil
}

func getGroupsHandler(c echo.Context) error {
	msg, err := n.Request("group.find", nil, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getGroupHandler(c echo.Context) error {
	query := fmt.Sprintf(`{"name": "%s"}`, c.Param("group"))
	msg, err := n.Request("group.get", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func createGroupHandler(c echo.Context) error {
	var g Group

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if g.Map(c) != nil {
		return ErrBadReqBody
	}

	data, err := json.Marshal(g)
	if err != nil {
		return ErrInternal
	}

	msg, err := n.Request("group.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateGroupHandler(c echo.Context) error {
	return ErrNotImplemented
}

func deleteGroupHandler(c echo.Context) error {
	return ErrNotImplemented
}
