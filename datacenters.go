/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// getDatacentersHandler : responds to GET /datacenters/ with a list of all
// datacenters
func getDatacentersHandler(c echo.Context) (err error) {
	var datacenters []Datacenter
	var body []byte
	var datacenter Datacenter

	au := authenticatedUser(c)
	if au.Admin == true {
		err = datacenter.FindAll(au, &datacenters)
	} else {
		datacenters, err = au.Datacenters()
	}

	if err != nil {
		return err
	}

	for i := 0; i < len(datacenters); i++ {
		datacenters[i].Improve()
	}

	if body, err = json.Marshal(datacenters); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// getDatacenterHandler : responds to GET /datacenter/:id:/ with the specified
// datacenter details
func getDatacenterHandler(c echo.Context) (err error) {
	var d Datacenter
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

// createDatacenterHandler : responds to POST /datacenters/ by creating a
// datacenter on the data store
func createDatacenterHandler(c echo.Context) (err error) {
	var d Datacenter
	var existing Datacenter
	var body []byte

	au := authenticatedUser(c)

	if au.GroupID == 0 {
		return c.JSONBlob(401, []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action"))
	}

	if d.Map(c) != nil {
		return ErrBadReqBody
	}

	d.GroupID = au.GroupID

	if err := existing.FindByName(d.Name, &existing); err == nil {
		return echo.NewHTTPError(409, "Specified datacenter already exists")
	}

	d.Save()

	if body, err = json.Marshal(d); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, body)
}

// updateDatacenterHandler : responds to PUT /datacenters/:id: by updating
// an existing datacenter
func updateDatacenterHandler(c echo.Context) (err error) {
	var d Datacenter
	var existing Datacenter
	var body []byte

	if d.Map(c) != nil {
		return ErrBadReqBody
	}

	au := authenticatedUser(c)
	if au.Admin != true {
		return ErrUnauthorized
	}

	if err := existing.FindByName(d.Name, &existing); err != nil {
		return echo.NewHTTPError(404, "Specified datacenter does not exists")
	}

	d.Save()

	if body, err = json.Marshal(d); err != nil {
		return ErrInternal
	}

	return c.JSONBlob(http.StatusOK, body)
}

// deleteDatacenterHandler : responds to DELETE /datacenters/:id: by deleting an
// existing datacenter
func deleteDatacenterHandler(c echo.Context) error {
	var d Datacenter

	au := authenticatedUser(c)

	id, err := strconv.Atoi(c.Param("datacenter"))
	if err = d.FindByID(id); err != nil {
		return err
	}

	if au.GroupID != d.GroupID {
		return ErrUnauthorized
	}

	if err := d.Delete(); err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}
