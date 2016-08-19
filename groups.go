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
	"strconv"
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

// getGroupsHandler : get all datacenters
func getGroupsHandler(c echo.Context) error {
	var msg *nats.Msg
	var err error

	au := authenticatedUser(c)
	if au.Admin == false {
		if body, err := getGroupByID(au.GroupID); err != nil {
			return err
		} else {
			body := []byte("[" + string(body) + "]")
			return c.JSONBlob(http.StatusOK, body)
		}
	}

	msg, err = n.Request("group.find", nil, 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

// getGroupHandler : get a group by id
func getGroupHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("group"))
	if body, err := getGroupByID(id); err != nil {
		return err
	} else {
		return c.JSONBlob(http.StatusOK, body)
	}
}

func getGroupByID(id int) (body []byte, err error) {
	query := fmt.Sprintf(`{"id": %d}`, id)
	msg, err := n.Request("group.get", []byte(query), 1*time.Second)
	body = msg.Data

	if err != nil {
		return body, ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return body, re.HTTPError
	}

	return body, nil
}

// createGroupHandler : Endpoint to create a datacenter
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

	return c.JSONBlob(http.StatusOK, msg.Data)
}

// updateGroupHandler : Updates a group trough its store
func updateGroupHandler(c echo.Context) error {
	var g Group
	if g.Map(c) != nil {
		return ErrBadReqBody
	}

	au := authenticatedUser(c)
	if au.Admin != true {
		return ErrUnauthorized
	}

	data, err := json.Marshal(g)
	if err != nil {
		return ErrInternal
	}

	msg, err := n.Request("group.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

// deleteGroupHandler : Deletes a group though its store
func deleteGroupHandler(c echo.Context) error {
	au := authenticatedUser(c)

	if au.Admin != true {
		return ErrUnauthorized
	}

	// Check if there is users on the group
	var users []User
	query := fmt.Sprintf(`{"group_id": %s}`, c.Param("group"))
	msg, err := n.Request("user.find", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}
	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}
	err = json.Unmarshal(msg.Data, &users)
	if err != nil {
		return ErrInternal
	}

	if len(users) > 0 {
		return ErrInternal
	}

	// Check if there is datacenters on the group
	var datacenters []Datacenter
	query = fmt.Sprintf(`{"group_id": %s}`, c.Param("group"))
	msg, err = n.Request("datacenter.find", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}
	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}
	err = json.Unmarshal(msg.Data, &datacenters)
	if err != nil {
		return ErrInternal
	}

	if len(datacenters) > 0 {
		return ErrInternal
	}

	query = fmt.Sprintf(`{"id": %s}`, c.Param("group"))
	msg, err = n.Request("group.del", []byte(query), 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.String(http.StatusOK, "")
}

// deleteUserFromGroupHandler : Deletes an user from a group
func deleteUserFromGroupHandler(c echo.Context) error {
	au := authenticatedUser(c)

	if au.Admin != true {
		return ErrUnauthorized
	}

	userid, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		return ErrBadReqBody
	}
	udata, err := getUser(userid)
	if err != nil {
		return ErrBadReqBody
	}

	var user User
	err = json.Unmarshal(udata, &user)
	if err != nil {
		return ErrBadReqBody
	}
	user.GroupID = 0
	user.Password = ""
	user.Salt = ""

	data, err := json.Marshal(user)
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

	return c.JSONBlob(http.StatusOK, []byte(""))

}

// addUserToGroupHandler : Adds an user to a group
func addUserToGroupHandler(c echo.Context) error {
	au := authenticatedUser(c)

	if au.Admin != true {
		return ErrUnauthorized
	}

	groupid, err := strconv.Atoi(c.Param("group"))
	if err != nil {
		return ErrBadReqBody
	}
	_, err = getGroup(groupid)
	if err != nil {
		return ErrBadReqBody
	}

	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	var payload struct {
		GroupID string `json:"groupid"`
		UserID  string `json:"userid"`
	}
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return ErrBadReqBody
	}

	userid, err := strconv.Atoi(payload.UserID)
	if err != nil {
		return ErrBadReqBody
	}
	udata, err := getUser(userid)
	if err != nil {
		return ErrBadReqBody
	}

	var user User
	err = json.Unmarshal(udata, &user)
	if err != nil {
		return ErrBadReqBody
	}
	user.GroupID = groupid
	user.Password = ""
	user.Salt = ""

	data, err = json.Marshal(user)
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

	return c.JSONBlob(http.StatusOK, []byte(""))
}

// addDatacenterToGroupHandler : Adds a datacenter to a group
func addDatacenterToGroupHandler(c echo.Context) error {
	au := authenticatedUser(c)

	if au.Admin != true {
		return ErrUnauthorized
	}

	groupid, err := strconv.Atoi(c.Param("group"))
	if err != nil {
		return ErrBadReqBody
	}
	_, err = getGroup(groupid)
	if err != nil {
		return ErrBadReqBody
	}

	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	var payload struct {
		GroupID      string `json:"groupid"`
		DatacenterID string `json:"datacenterid"`
	}
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return ErrBadReqBody
	}

	datacenterid, err := strconv.Atoi(payload.DatacenterID)
	if err != nil {
		return ErrBadReqBody
	}
	ddata, err := getDatacenterByID(datacenterid)
	if err != nil {
		return ErrBadReqBody
	}

	var datacenter Datacenter
	err = json.Unmarshal(ddata, &datacenter)
	if err != nil {
		return ErrBadReqBody
	}
	datacenter.GroupID = groupid

	data, err = json.Marshal(datacenter)
	if err != nil {
		return ErrBadReqBody
	}

	msg, err := n.Request("datacenter.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, []byte(""))
}

// deleteDatacenterFromGroupHandler : Deletes a datacenter from a group
func deleteDatacenterFromGroupHandler(c echo.Context) error {
	au := authenticatedUser(c)

	if au.Admin != true {
		return ErrUnauthorized
	}

	groupid, err := strconv.Atoi(c.Param("group"))
	if err != nil {
		return ErrBadReqBody
	}
	_, err = getGroup(groupid)
	if err != nil {
		return ErrBadReqBody
	}

	datacenterid, err := strconv.Atoi(c.Param("datacenter"))
	if err != nil {
		return ErrBadReqBody
	}
	_, err = getGroup(groupid)
	if err != nil {
		return ErrBadReqBody
	}

	ddata, err := getDatacenterByID(datacenterid)
	if err != nil {
		return ErrBadReqBody
	}

	var datacenter Datacenter
	err = json.Unmarshal(ddata, &datacenter)
	if err != nil {
		return ErrBadReqBody
	}
	datacenter.GroupID = 0

	data, err := json.Marshal(datacenter)
	if err != nil {
		return ErrBadReqBody
	}

	msg, err := n.Request("datacenter.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusOK, []byte(""))
}
