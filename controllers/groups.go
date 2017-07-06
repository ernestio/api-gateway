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
	return genericList(c, "group", groups.List)
}

// GetGroupHandler : responds to GET /groups/:id:/ with the specified
// group details
func GetGroupHandler(c echo.Context) (err error) {
	return genericGet(c, "group", groups.Get)
}

// CreateGroupHandler : responds to POST /groups/ by creating a group
// on the data store
func CreateGroupHandler(c echo.Context) (err error) {
	return genericCreate(c, "group", groups.Create)
}

// UpdateGroupHandler : responds to PUT /groups/:id: by updating an existing
// group
func UpdateGroupHandler(c echo.Context) (err error) {
	return genericUpdate(c, "group", groups.Update)
}

// DeleteGroupHandler : responds to DELETE /groups/:id: by deleting an
// existing group
func DeleteGroupHandler(c echo.Context) (err error) {
	return genericDelete(c, "group", groups.Delete)
}

// DeleteUserFromGroupHandler : Deletes an user from a group
func DeleteUserFromGroupHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "groups/rm_user")
	if st == 200 {
		u := c.Param("user")
		st, b = groups.RmUser(au, u)
	}

	return h.Respond(c, st, b)
}

// AddUserToGroupHandler : Adds an user to a group
func AddUserToGroupHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "groups/add_user")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	g := c.Param("group")

	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = groups.AddUser(au, g, body)
	}

	return h.Respond(c, st, b)
}

// AddDatacenterToGroupHandler : Adds a datacenter to a group
func AddDatacenterToGroupHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "groups/add_datacenter")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	g := c.Param("group")

	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = groups.AddDatacenter(au, g, body)
	}

	return h.Respond(c, st, b)
}

// DeleteDatacenterFromGroupHandler : Deletes a datacenter from a group
func DeleteDatacenterFromGroupHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "groups/rm_datacenter")
	if st == 200 {
		g := c.Param("group")
		d := c.Param("datacenter")
		st, b = groups.RmDatacenter(au, g, d)
	}

	return h.Respond(c, st, b)
}
