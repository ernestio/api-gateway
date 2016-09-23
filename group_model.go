/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/nats-io/nats"
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

func (g *Group) findByName(name string) (err error) {
	var msg *nats.Msg

	query := `{"name": "` + name + `"}`
	if msg, err = n.Request("group.get", []byte(query), 1*time.Second); err != nil {
		return ErrGatewayTimeout
	}
	if strings.Contains(string(msg.Data), `"error"`) {
		return errors.New(`"Specified group does not exist"`)
	}
	if err = json.Unmarshal(msg.Data, &g); err != nil {
		return errors.New(`"Specified group does not exist"`)
	}

	return nil
}
