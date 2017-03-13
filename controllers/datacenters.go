/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetDatacentersHandler : responds to GET /datacenters/ with a list of all
// datacenters
func GetDatacentersHandler(c echo.Context) (err error) {
	var datacenters []models.Datacenter
	var body []byte
	var datacenter models.Datacenter

	au := AuthenticatedUser(c)
	if au.Admin == true {
		err = datacenter.FindAll(au, &datacenters)
	} else {
		datacenters, err = au.Datacenters()
	}

	if err != nil {
		return err
	}

	for i := 0; i < len(datacenters); i++ {
		datacenters[i].Redact()
		datacenters[i].Improve()
	}

	if body, err = json.Marshal(datacenters); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// GetDatacenterHandler : responds to GET /datacenter/:id:/ with the specified
// datacenter details
func GetDatacenterHandler(c echo.Context) (err error) {
	var d models.Datacenter
	var body []byte

	id, _ := strconv.Atoi(c.Param("datacenter"))
	if err := d.FindByID(id); err != nil {
		return err
	}

	if body, err = json.Marshal(d); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, body)
}

// CreateDatacenterHandler : responds to POST /datacenters/ by creating a
// datacenter on the data store
func CreateDatacenterHandler(c echo.Context) (err error) {
	var d models.Datacenter
	var existing models.Datacenter
	var body []byte

	au := AuthenticatedUser(c)

	if au.GroupID == 0 {
		return c.JSONBlob(401, []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action"))
	}

	if d.Map(c) != nil {
		return h.ErrBadReqBody
	}

	err = d.Validate()
	if err != nil {
		return h.ErrBadReqBody
	}

	d.GroupID = au.GroupID

	if err := existing.FindByName(d.Name, &existing); err == nil {
		return echo.NewHTTPError(409, "Specified datacenter already exists")
	}

	if err = d.Save(); err != nil {
		log.Println(err)
	}

	if body, err = json.Marshal(d); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, body)
}

// UpdateDatacenterHandler : responds to PUT /datacenters/:id: by updating
// an existing datacenter
func UpdateDatacenterHandler(c echo.Context) (err error) {
	var d models.Datacenter
	var existing models.Datacenter
	var body []byte

	if d.Map(c) != nil {
		return h.ErrBadReqBody
	}

	au := AuthenticatedUser(c)

	id, err := strconv.Atoi(c.Param("datacenter"))
	if err = existing.FindByID(id); err != nil {
		return err
	}

	if au.GroupID != au.GroupID {
		return h.ErrUnauthorized
	}

	existing.Username = d.Username
	existing.Password = d.Password
	existing.AccessKeyID = d.AccessKeyID
	existing.SecretAccessKey = d.SecretAccessKey

	if err = existing.Save(); err != nil {
		log.Println(err)
	}

	if body, err = json.Marshal(d); err != nil {
		return h.ErrInternal
	}

	return c.JSONBlob(http.StatusOK, body)
}

// DeleteDatacenterHandler : responds to DELETE /datacenters/:id: by deleting an
// existing datacenter
func DeleteDatacenterHandler(c echo.Context) error {
	var d models.Datacenter

	au := AuthenticatedUser(c)

	id, err := strconv.Atoi(c.Param("datacenter"))
	if err = d.FindByID(id); err != nil {
		return err
	}

	if au.GroupID != d.GroupID {
		return h.ErrUnauthorized
	}

	ss, err := d.Services()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	if len(ss) > 0 {
		return echo.NewHTTPError(400, "Existing services are referring to this datacenter.")
	}

	if err := d.Delete(); err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}
