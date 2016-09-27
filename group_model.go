/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

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

func (g *Group) FindByName(name string, group *Group) (err error) {
	var res []byte

	query := `{"name": "` + name + `"}`
	if res, err = Query("group.get", query); err != nil {
		return err
	}
	if strings.Contains(string(res), `"error"`) {
		return errors.New(`"Specified group does not exist"`)
	}
	if err = json.Unmarshal(res, &group); err != nil {
		return errors.New(`"Specified group does not exist"`)
	}

	return nil
}

// FindAll : Searches for all groups on the store current user
// has access to
func (g *Group) FindAll(au User, groups *[]Group) (err error) {
	var query string
	var res []byte

	if !au.Admin {
		query = fmt.Sprintf(`{"group_id": %d}`, au.GroupID)
	}

	if res, err = Query("group.find", query); err != nil {
		return err
	}

	err = json.Unmarshal(res, &groups)
	if err != nil {
		return ErrInternal
	}

	return nil
}

func (g *Group) FindByID(id int) (err error) {
	var res []byte
	query := fmt.Sprintf(`{"id": %d}`, id)
	if res, err = Query("group.get", query); err != nil {
		return err
	}
	json.Unmarshal(res, &g)
	return nil
}

// Save : calls group.set with the marshalled current group
func (g *Group) Save() (err error) {
	var res []byte

	data, err := json.Marshal(g)
	if err != nil {
		return ErrBadReqBody
	}

	if res, err = Query("group.set", string(data)); err != nil {
		return err
	}

	if err := json.Unmarshal(res, &g); err != nil {
		return ErrInternal
	}

	return nil
}

// Delete : will delete a group by its id
func (g *Group) Delete() (err error) {
	query := fmt.Sprintf(`{"id": %s}`, g.ID)
	if _, err := Query("group.del", query); err != nil {
		return err
	}

	return nil
}

// Users : Get the users related with current group
func (g *Group) Users() (users []User, err error) {
	var u User
	u.GroupID = g.ID
	err = u.FindAll(&users)

	return users, err
}

// Datacenters : Get the datacenters related with current group
func (g *Group) Datacenters() (datacenters []Datacenter, err error) {
	var d Datacenter
	err = d.FindByGroupID(g.ID, &datacenters)

	return datacenters, err
}
