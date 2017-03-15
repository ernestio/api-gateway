/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"

	"github.com/Sirupsen/logrus"
	h "github.com/ernestio/api-gateway/helpers"
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
func (g *Group) Map(data []byte) error {
	if err := json.Unmarshal(data, &g); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	if err := g.Validate(); err != nil {
		h.L.WithFields(logrus.Fields{
			"input":         string(data),
			"error_message": err.Error(),
		}).Error("Invalid input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// FindByName : Searches for all groups with a name equal to the specified
func (g *Group) FindByName(name string, group *Group) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel("group").GetBy(query, group); err != nil {
		return err
	}
	return nil
}

// FindAll : Searches for all groups on the store current user
// has access to
func (g *Group) FindAll(au User, groups *[]Group) (err error) {
	query := make(map[string]interface{})
	if !au.Admin {
		query["group_id"] = au.GroupID
	}
	if err := NewBaseModel("group").FindBy(query, groups); err != nil {
		return err
	}
	return nil
}

// FindByID : Gets a model by its id
func (g *Group) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	return NewBaseModel("group").GetBy(query, g)
}

// Save : calls group.set with the marshalled current group
func (g *Group) Save() (err error) {
	if err := NewBaseModel("group").Save(g); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a group by its id
func (g *Group) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = g.ID
	if err := NewBaseModel("group").Delete(query); err != nil {
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
