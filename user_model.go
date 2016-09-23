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

// User holds the user response from user-store
type User struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldpassword,omitempty"`
	Salt        string `json:"salt,omitempty"`
	Admin       bool   `json:"admin"`
}

// Validate vaildate all of the user's input
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("User username is empty")
	}

	if u.Password == "" {
		return errors.New("User password is empty")
	}

	return nil
}

// Map a user from a request's body and validates the input
func (u *User) Map(c echo.Context) *echo.HTTPError {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &u)
	if err != nil {
		return ErrBadReqBody
	}

	err = u.Validate()
	if err != nil {
		return ErrBadReqBody
	}

	return nil
}

func (u *User) findByUserName(name string) (err error) {
	var msg *nats.Msg

	query := `{"username": "` + name + `"}`
	if msg, err = n.Request("user.get", []byte(query), 1*time.Second); err != nil {
		return ErrGatewayTimeout
	}
	if strings.Contains(string(msg.Data), `"error"`) {
		return errors.New(`"Specified username does not exist"`)
	}
	if err = json.Unmarshal(msg.Data, &u); err != nil {
		return errors.New(`"Specified username does not exist"`)
	}

	return nil
}

func (u *User) save() (err error) {
	data, err := json.Marshal(u)
	if err != nil {
		return ErrBadReqBody
	}

	msg, err := n.Request("user.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return nil
}
