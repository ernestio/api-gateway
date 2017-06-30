/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
	"github.com/labstack/echo"
)

// GetServiceBuildHandler : gets the details of a specific service build
func GetServiceBuildHandler(c echo.Context) (err error) {
	var list []views.ServiceRender

	au := AuthenticatedUser(c)
	query := h.GetParamFilter(c)
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if list, err = getServicesOutput(query); err != nil {
		return c.JSONBlob(500, []byte(err.Error()))
	}

	if len(list) > 0 {
		return c.JSON(http.StatusOK, list[0])
	}
	return c.JSON(http.StatusNotFound, nil)
}

// DelServiceBuildHandler : will delete the specified build from a service
func DelServiceBuildHandler(c echo.Context) (err error) {
	var services []models.Service
	var s models.Service

	au := AuthenticatedUser(c)
	query := make(map[string]interface{})
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}
	id := c.Param("build")
	query["id"] = id
	if id == "" {
		h.L.Debug("Empty id")
		return c.JSONBlob(400, []byte("Invalid build id"))
	}

	if err := s.Find(query, &services); err != nil {
		return c.JSONBlob(404, []byte("Not found"))
	}

	if len(services) == 0 {
		h.L.Debug("Build " + id + " not found")
		return c.JSONBlob(404, []byte("Not found"))
	}
	if err := services[0].Delete(); err != nil {
		h.L.Warning(err.Error())
		return c.JSONBlob(500, []byte("Oops something went wrong"))
	}

	return c.JSONBlob(501, []byte("Not implemented"))
}
