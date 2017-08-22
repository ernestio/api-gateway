/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/projects"
	"github.com/labstack/echo"
)

// GetDatacentersHandler : responds to GET /projects/ with a list of all
// projects
func GetDatacentersHandler(c echo.Context) (err error) {
	return genericList(c, "project", projects.List)
}

// GetDatacenterHandler : responds to GET /projects/:id:/ with the specified
// projects details
func GetDatacenterHandler(c echo.Context) (err error) {
	return genericGet(c, "project", projects.Get)
}

// CreateDatacenterHandler : responds to POST /projects/ by creating a
// project on the data store
func CreateDatacenterHandler(c echo.Context) (err error) {
	return genericCreate(c, "project", projects.Create)
}

// UpdateDatacenterHandler : responds to PUT /projects/:id: by updating
// an existing project
func UpdateDatacenterHandler(c echo.Context) (err error) {
	return genericUpdate(c, "project", projects.Update)
}

// DeleteDatacenterHandler : responds to DELETE /projects/:id: by deleting an
// existing project
func DeleteDatacenterHandler(c echo.Context) error {
	return genericDelete(c, "project", projects.Delete)
}
