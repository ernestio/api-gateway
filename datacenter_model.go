/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/labstack/echo"
)

// Datacenter holds the datacenter response from datacenter-store
type Datacenter struct {
	ID              int    `json:"id"`
	GroupID         int    `json:"group_id"`
	GroupName       string `json:"group_name"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Region          string `json:"region"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	VCloudURL       string `json:"vcloud_url"`
	VseURL          string `json:"vse_url"`
	ExternalNetwork string `json:"external_network"`
	AccessKeyID     string `json:"aws_access_key_id,omitempty"`
	SecretAccessKey string `json:"aws_secret_access_key,omitempty"`
}

// Validate the datacenter
func (d *Datacenter) Validate() error {
	if d.Name == "" {
		return errors.New("Datacenter name is empty")
	}

	if d.Type == "" {
		return errors.New("Datacenter type is empty")
	}

	if d.Username == "" {
		return errors.New("Datacenter username is empty")
	}

	if d.Type == "vcloud" && d.VCloudURL == "" {
		return errors.New("Datacenter vcloud url is empty")
	}

	return nil
}

// Map : maps a datacenter from a request's body and validates the input
func (d *Datacenter) Map(c echo.Context) *echo.HTTPError {
	body := c.Request().Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return ErrBadReqBody
	}

	return nil
}

// FindByName : Searches for all datacenters with a name equal to the specified
func (d *Datacenter) FindByName(name string, datacenter *Datacenter) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel("datacenter").GetBy(query, datacenter); err != nil {
		return err
	}
	return nil
}

// FindByGroupID : Searches for all datacenters on the store current user
// has access to with the specified group id
func (d *Datacenter) FindByGroupID(id int, datacenters *[]Datacenter) (err error) {
	query := make(map[string]interface{})
	query["group_id"] = id
	if err := NewBaseModel("datacenter").FindBy(query, datacenters); err != nil {
		return err
	}
	return nil
}

// FindByNameAndGroupID : Searches for all datacenters with a name equal to the specified
func (d *Datacenter) FindByNameAndGroupID(name string, id int, datacenters *[]Datacenter) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	query["group_id"] = id
	if err := NewBaseModel("datacenter").FindBy(query, datacenters); err != nil {
		return err
	}
	return nil
}

// FindByID : Gets a model by its id
func (d *Datacenter) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	if err := NewBaseModel("datacenter").GetBy(query, d); err != nil {
		return err
	}
	return nil
}

// FindAll : Searches for all groups on the store current user
// has access to
func (d *Datacenter) FindAll(au User, datacenters *[]Datacenter) (err error) {
	query := make(map[string]interface{})
	if err := NewBaseModel("datacenter").FindBy(query, datacenters); err != nil {
		return err
	}
	return nil
}

// Save : calls datacenter.set with the marshalled current group
func (d *Datacenter) Save() (err error) {
	if err := NewBaseModel("datacenter").Save(d); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a datacenter by its id
func (d *Datacenter) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = d.ID
	if err := NewBaseModel("datacenter").Delete(query); err != nil {
		return err
	}
	return nil
}

// Redact : removes all sensitive fields from the return
// data before outputting to the user
func (d *Datacenter) Redact() {
	d.AccessKeyID = ""
	d.SecretAccessKey = ""
}

// Improve : adds extra data as group name
func (d *Datacenter) Improve() {
	g := d.Group()
	d.GroupName = g.Name
}

// Group : Gets the related datacenter group if any
func (d *Datacenter) Group() (group Group) {
	if err := group.FindByID(d.GroupID); err != nil {
		log.Println(err)
	}

	return group
}

// Services : Get the services related with current datacenter
func (d *Datacenter) Services() (services []Service, err error) {
	var s Service
	err = s.FindByDatacenterID(d.ID, &services)

	return services, err
}
