/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/datacenters"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetDatacentersHandler : responds to GET /datacenters/ with a list of all
// datacenters
func GetDatacentersHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := datacenters.List(au)

	return c.JSONBlob(st, b)
}

// GetDatacenterHandler : responds to GET /datacenter/:id:/ with the specified
// datacenter details
func GetDatacenterHandler(c echo.Context) (err error) {
	d := c.Param("datacenter")
	st, b := datacenters.Get(d)

	return c.JSONBlob(st, b)
}

// CreateDatacenterHandler : responds to POST /datacenters/ by creating a
// datacenter on the data store
func CreateDatacenterHandler(c echo.Context) (err error) {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = datacenters.Create(au, body)
	}

	return c.JSONBlob(s, b)
}

// UpdateDatacenterHandler : responds to PUT /datacenters/:id: by updating
// an existing datacenter
func UpdateDatacenterHandler(c echo.Context) (err error) {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)
	d := c.Param("datacenter")

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = datacenters.Update(au, d, body)
	}

	return c.JSONBlob(s, b)
}

// DeleteDatacenterHandler : responds to DELETE /datacenters/:id: by deleting an
// existing datacenter
func DeleteDatacenterHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	d := c.Param("datacenter")

	s, b := datacenters.Delete(au, d)

	return c.JSONBlob(s, b)
}
