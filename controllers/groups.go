/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetGroupsHandler : responds to GET /groups/ with a list of all
// groups
func GetGroupsHandler(c echo.Context) (err error) {
	var groups []models.Group
	var body []byte
	var group models.Group

	au := AuthenticatedUser(c)
	if au.Admin == true {
		if err := group.FindAll(au, &groups); err != nil {
			log.Println(err)
		}
	} else {
		if err := group.FindByID(au.GroupID); err != nil {
			log.Println(err)
		}
		groups = append(groups, group)
	}

	if body, err = json.Marshal(groups); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// GetGroupHandler : responds to GET /groups/:id:/ with the specified
// group details
func GetGroupHandler(c echo.Context) (err error) {
	var g models.Group
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

// CreateGroupHandler : responds to POST /groups/ by creating a group
// on the data store
func CreateGroupHandler(c echo.Context) (err error) {
	var g models.Group
	var existing models.Group
	var body []byte

	if AuthenticatedUser(c).Admin != true {
		return h.ErrUnauthorized
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}

	if g.Map(data) != nil {
		return h.ErrBadReqBody
	}

	if err := existing.FindByName(g.Name, &existing); err == nil {
		return echo.NewHTTPError(409, "Specified group already exists")
	}

	if err = g.Save(); err != nil {
		log.Println(err)
	}

	if body, err = json.Marshal(g); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, body)
}

// UpdateGroupHandler : responds to PUT /groups/:id: by updating an existing
// group
func UpdateGroupHandler(c echo.Context) (err error) {
	var g models.Group
	var existing models.Group
	var body []byte

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}

	if g.Map(data) != nil {
		return h.ErrBadReqBody
	}

	au := AuthenticatedUser(c)
	if au.Admin != true {
		return h.ErrUnauthorized
	}

	if err := existing.FindByName(g.Name, &existing); err != nil {
		return echo.NewHTTPError(404, "Specified group does not exists")
	}

	if err = g.Save(); err != nil {
		log.Println(err)
	}

	if body, err = json.Marshal(g); err != nil {
		return h.ErrInternal
	}

	return c.JSONBlob(http.StatusOK, body)
}

// DeleteGroupHandler : responds to DELETE /groups/:id: by deleting an
// existing group
func DeleteGroupHandler(c echo.Context) (err error) {
	var g models.Group
	var users []models.User
	var datacenters []models.Datacenter

	au := AuthenticatedUser(c)

	if au.Admin != true {
		return h.ErrUnauthorized
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

// DeleteUserFromGroupHandler : Deletes an user from a group
func DeleteUserFromGroupHandler(c echo.Context) error {
	var user models.User
	au := AuthenticatedUser(c)

	if au.Admin == false {
		return h.ErrUnauthorized
	}

	if err := user.FindByID(c.Param("user"), &user); err != nil {
		log.Println(err)
	}
	user.GroupID = 0
	user.Password = ""
	user.Salt = ""
	if err := user.Save(); err != nil {
		return h.ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, []byte("User "+user.Username+" successfully removed from group"))
}

// AddUserToGroupHandler : Adds an user to a group
func AddUserToGroupHandler(c echo.Context) error {
	var group models.Group
	var user models.User
	var payload map[string]string

	au := AuthenticatedUser(c)

	if au.Admin != true {
		return h.ErrUnauthorized
	}

	if err := group.FindByName(c.Param("group"), &group); err != nil {
		return h.ErrBadReqBody
	}

	body := c.Request().Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return h.ErrBadReqBody
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		return h.ErrBadReqBody
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

// AddDatacenterToGroupHandler : Adds a datacenter to a group
func AddDatacenterToGroupHandler(c echo.Context) error {
	var group models.Group
	var datacenter models.Datacenter
	var payload map[string]string

	au := AuthenticatedUser(c)
	if au.Admin != true {
		return h.ErrUnauthorized
	}

	groupID, err := strconv.Atoi(c.Param("group"))
	if err != nil {
		return h.ErrBadReqBody
	}

	body := c.Request().Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return h.ErrBadReqBody
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		return h.ErrBadReqBody
	}

	datacenterID, err := strconv.Atoi(payload["datacenterid"])
	if err != nil {
		return h.ErrBadReqBody
	}

	if err := group.FindByID(groupID); err != nil {
		return h.ErrBadReqBody
	}

	if err := datacenter.FindByID(datacenterID); err != nil {
		return h.ErrBadReqBody
	}

	datacenter.GroupID = groupID
	if err = datacenter.Save(); err != nil {
		log.Println(err)
	}

	return c.JSONBlob(http.StatusOK, []byte("Datacenter successfully added to group "+group.Name))
}

// DeleteDatacenterFromGroupHandler : Deletes a datacenter from a group
func DeleteDatacenterFromGroupHandler(c echo.Context) error {
	var group models.Group
	var datacenter models.Datacenter

	au := AuthenticatedUser(c)

	if au.Admin != true {
		return h.ErrUnauthorized
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
	if err = datacenter.Save(); err != nil {
		log.Println(err)
	}

	return c.JSONBlob(http.StatusOK, []byte("Datacenter successfully removed from group "+group.Name))
}
