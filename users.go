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

// User holds the user response from user-store
type User struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

// Validate the user
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("User username is empty")
	}

	if u.Password == "" {
		return errors.New("User password is empty")
	}

	if u.GroupID == "" {
		return errors.New("User group is empty")
	}

	return nil
}

// Map : maps a user from a request's body and validates the input
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

func getUsersHandler(c echo.Context) error {
	msg, err := n.Request("user.find", nil, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getUserHandler(c echo.Context) error {
	query := fmt.Sprintf(`{"username": "%s"}`, c.Param("user"))
	msg, err := n.Request("user.get", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func createUserHandler(c echo.Context) error {
	var u User

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if u.Map(c) != nil {
		return ErrBadReqBody
	}

	data, err := json.Marshal(u)
	if err != nil {
		return ErrInternal
	}

	msg, err := n.Request("user.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateUserHandler(c echo.Context) error {
	var u User
	var qu User

	if u.Map(c) != nil {
		return ErrBadReqBody
	}

	data, err := json.Marshal(u)
	if err != nil {
		return ErrInternal
	}

	// Check if authenticated user is admin or updating itself
	au := authenticatedUser(c)
	if au.Username != u.Username && au.Admin != true {
		return ErrUnauthorized
	}

	// Check user exists
	msg, err := n.Request("user.get", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	err = json.Unmarshal(msg.Data, &qu)
	if err != nil {
		return ErrInternal
	}

	if qu.ID == "" {
		return ErrNotFound
	}

	// update the user
	msg, err = n.Request("user.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func deleteUserHandler(c echo.Context) error {
	query := fmt.Sprintf(`{"username": "%s"}`, c.Param("user"))
	msg, err := n.Request("user.del", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.String(http.StatusOK, "")
}
