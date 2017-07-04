/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/users"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetUsersHandler : responds to GET /users/ with a list of all
// users for admin, and all users in your group for other
// users
func GetUsersHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := users.List(au)

	return c.JSONBlob(st, b)
}

// GetUserHandler : responds to GET /users/:id:/ with the specified
// user details
func GetUserHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	u := c.Param("user")
	st, b := users.Get(au, u)

	return c.JSONBlob(st, b)
}

// CreateUserHandler : responds to POST /users/ by creating a user
// on the data store
func CreateUserHandler(c echo.Context) error {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = users.Create(au, body)
	}

	return c.JSONBlob(s, b)
}

// UpdateUserHandler : responds to PUT /users/:id: by updating an existing
// user
func UpdateUserHandler(c echo.Context) error {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)
	d := c.Param("user")

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = users.Update(au, d, body)
	}

	return c.JSONBlob(s, b)
}

// DeleteUserHandler : responds to DELETE /users/:id: by deleting an
// existing user
func DeleteUserHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	u := c.Param("user")
	st, b := users.Delete(au, u)

	return c.JSONBlob(st, b)
}
