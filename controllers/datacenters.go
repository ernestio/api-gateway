/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/datacenters"
	"github.com/labstack/echo"
)

// GetDatacentersHandler : responds to GET /datacenters/ with a list of all
// datacenters
func GetDatacentersHandler(c echo.Context) (err error) {
	return genericList(c, "project", datacenters.List)
}

// GetDatacenterHandler : responds to GET /datacenter/:id:/ with the specified
// datacenter details
func GetDatacenterHandler(c echo.Context) (err error) {
	return genericGet(c, "project", datacenters.Get)
}

// CreateDatacenterHandler : responds to POST /datacenters/ by creating a
// datacenter on the data store
func CreateDatacenterHandler(c echo.Context) (err error) {
	return genericCreate(c, "project", datacenters.Create)
}

// UpdateDatacenterHandler : responds to PUT /datacenters/:id: by updating
// an existing datacenter
func UpdateDatacenterHandler(c echo.Context) (err error) {
	return genericUpdate(c, "project", datacenters.Update)
}

// DeleteDatacenterHandler : responds to DELETE /datacenters/:id: by deleting an
// existing datacenter
func DeleteDatacenterHandler(c echo.Context) error {
	return genericDelete(c, "project", datacenters.Delete)
}
