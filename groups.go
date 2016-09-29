/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// getGroupsHandler : responds to GET /groups/ with a list of all
// groups
func getGroupsHandler(c echo.Context) (err error) {
	var groups []Group
	var body []byte
	var group Group

	au := authenticatedUser(c)
	if au.Admin == true {
		group.FindAll(au, &groups)
	} else {
		group.FindByID(au.GroupID)
		groups = append(groups, group)
	}

	if body, err = json.Marshal(groups); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// getGroupHandler : responds to GET /groups/:id:/ with the specified
// group details
func getGroupHandler(c echo.Context) (err error) {
	var g Group
	var body []byte

	id, _ := strconv.Atoi(c.Param("group"))
	if err := g.FindByID(id); err != nil {
		return err
	}

	if body, err = json.Marshal(g); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, body)
}

// createGroupHandler : responds to POST /groups/ by creating a group
// on the data store
func createGroupHandler(c echo.Context) (err error) {
	var g Group
	var existing Group
	var body []byte

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if g.Map(c) != nil {
		return ErrBadReqBody
	}

	if err := existing.FindByName(g.Name, &existing); err == nil {
		return echo.NewHTTPError(409, "Specified group already exists")
	}

	g.Save()

	if body, err = json.Marshal(g); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, body)
}

// updateGroupHandler : responds to PUT /groups/:id: by updating an existing
// group
func updateGroupHandler(c echo.Context) (err error) {
	var g Group
	var existing Group
	var body []byte

	if g.Map(c) != nil {
		return ErrBadReqBody
	}

	au := authenticatedUser(c)
	if au.Admin != true {
		return ErrUnauthorized
	}

	if err := existing.FindByName(g.Name, &existing); err != nil {
		return echo.NewHTTPError(404, "Specified group does not exists")
	}

	g.Save()

	if body, err = json.Marshal(g); err != nil {
		return ErrInternal
	}

	return c.JSONBlob(http.StatusOK, body)
}

// deleteGroupHandler : responds to DELETE /groups/:id: by deleting an
// existing group
func deleteGroupHandler(c echo.Context) (err error) {
	var g Group
	var users []User
	var datacenters []Datacenter

	au := authenticatedUser(c)

	if au.Admin != true {
		return ErrUnauthorized
	}

	id, err := strconv.Atoi(c.Param("group"))
	if err = g.FindByID(id); err != nil {
		return err
	}

	// Check if there are any users on the group
	if users, err = g.Users(); err != nil {
		return err
	}

	if len(users) > 0 {
		return echo.NewHTTPError(400, "This group has users assigned to it, please remove the users before performing this action")
	}

	// Check if there are any datacenters on the group
	if datacenters, err = g.Datacenters(); err != nil {
		return err
	}

	if len(datacenters) > 0 {
		return echo.NewHTTPError(400, "This group has datacenters assigned to it, please remove the datacenters before performing this action")
	}

	if err := g.Delete(); err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}

// deleteUserFromGroupHandler : Deletes an user from a group
func deleteUserFromGroupHandler(c echo.Context) error {
	var user User
	au := authenticatedUser(c)

	if au.Admin == false {
		return ErrUnauthorized
	}

	user.FindByID(c.Param("user"), &user)
	user.GroupID = 0
	user.Password = ""
	user.Salt = ""
	if err := user.Save(); err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, []byte("User "+user.Username+" successfully removed from group"))
}

// addUserToGroupHandler : Adds an user to a group
func addUserToGroupHandler(c echo.Context) error {
	var group Group
	var user User
	var payload map[string]string

	au := authenticatedUser(c)

	if au.Admin != true {
		return ErrUnauthorized
	}

	if err := group.FindByName(c.Param("group"), &group); err != nil {
		return ErrBadReqBody
	}

	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		return ErrBadReqBody
	}

	if err := user.FindByUserName(payload["username"], &user); err != nil {
		return err
	}

	user.GroupID = group.ID
	user.Password = ""
	user.Salt = ""
	if err := user.Save(); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, []byte("User "+user.Username+" successfully added to group "+group.Name))
}

// addDatacenterToGroupHandler : Adds a datacenter to a group
func addDatacenterToGroupHandler(c echo.Context) error {
	var group Group
	var datacenter Datacenter
	var payload map[string]string

	au := authenticatedUser(c)
	if au.Admin != true {
		return ErrUnauthorized
	}

	groupID, err := strconv.Atoi(c.Param("group"))
	if err != nil {
		return ErrBadReqBody
	}

	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		return ErrBadReqBody
	}

	datacenterID, err := strconv.Atoi(payload["datacenterid"])
	if err != nil {
		return ErrBadReqBody
	}

	if err := group.FindByID(groupID); err != nil {
		return ErrBadReqBody
	}

	if err := datacenter.FindByID(datacenterID); err != nil {
		return ErrBadReqBody
	}

	datacenter.GroupID = groupID
	datacenter.Save()

	return c.JSONBlob(http.StatusOK, []byte("Datacenter successfully added to group "+group.Name))
}

// deleteDatacenterFromGroupHandler : Deletes a datacenter from a group
func deleteDatacenterFromGroupHandler(c echo.Context) error {
	var group Group
	var datacenter Datacenter

	au := authenticatedUser(c)

	if au.Admin != true {
		return ErrUnauthorized
	}

	groupid, err := strconv.Atoi(c.Param("group"))
	if err = group.FindByID(groupid); err != nil {
		return err
	}

	datacenterid, err := strconv.Atoi(c.Param("datacenter"))
	if err = datacenter.FindByID(datacenterid); err != nil {
		return err
	}

	datacenter.GroupID = 0
	datacenter.Save()

	return c.JSONBlob(http.StatusOK, []byte("Datacenter successfully removed from group "+group.Name))
}
