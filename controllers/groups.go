/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/groups"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetGroupsHandler : responds to GET /groups/ with a list of all
// groups
func GetGroupsHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := groups.List(au)

	return c.JSONBlob(st, b)
}

// GetGroupHandler : responds to GET /groups/:id:/ with the specified
// group details
func GetGroupHandler(c echo.Context) (err error) {
	g := c.Param("group")
	st, b := groups.Get(g)

	return c.JSONBlob(st, b)
}

// CreateGroupHandler : responds to POST /groups/ by creating a group
// on the data store
func CreateGroupHandler(c echo.Context) (err error) {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = groups.Create(au, body)
	}

	return c.JSONBlob(s, b)
}

// UpdateGroupHandler : responds to PUT /groups/:id: by updating an existing
// group
func UpdateGroupHandler(c echo.Context) (err error) {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = groups.Update(au, body)
	}

	return c.JSONBlob(s, b)
}

// DeleteGroupHandler : responds to DELETE /groups/:id: by deleting an
// existing group
func DeleteGroupHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	g := c.Param("group")
	s, b := groups.Delete(au, g)

	return c.JSONBlob(s, b)
}

// DeleteUserFromGroupHandler : Deletes an user from a group
func DeleteUserFromGroupHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	u := c.Param("user")
	s, b := groups.RmUser(au, u)

	return c.JSONBlob(s, b)
}

// AddUserToGroupHandler : Adds an user to a group
func AddUserToGroupHandler(c echo.Context) error {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)
	g := c.Param("group")

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = groups.AddUser(au, g, body)
	}

	return c.JSONBlob(s, b)
}

// AddDatacenterToGroupHandler : Adds a datacenter to a group
func AddDatacenterToGroupHandler(c echo.Context) error {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)
	g := c.Param("group")
	d := c.Param("datacenterid")

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = groups.AddDatacenter(au, g, d, body)
	}

	return c.JSONBlob(s, b)
}

// DeleteDatacenterFromGroupHandler : Deletes a datacenter from a group
func DeleteDatacenterFromGroupHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	g := c.Param("group")
	d := c.Param("datacenter")
	s, b := groups.RmDatacenter(au, g, d)

	return c.JSONBlob(s, b)
}
