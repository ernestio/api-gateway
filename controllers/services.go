/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"strings"

	"github.com/ernestio/api-gateway/controllers/services"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetServicesHandler : responds to GET /services/ with a list ahorized services
func GetServicesHandler(c echo.Context) (err error) {
	return genericList(c, "service", services.List)
}

// GetServiceBuildsHandler : gets the list of builds for the specified
// service
func GetServiceBuildsHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/builds")
	if st == 200 {
		p := h.GetParamFilter(c)
		st, b = services.Builds(au, p)
	}

	return h.Respond(c, st, b)
}

// GetServiceHandler : responds to GET /services/:service with the
// details of an existing service
func GetServiceHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/get")
	if st == 200 {
		p := h.GetParamFilter(c)
		st, b = services.Get(au, p)
	}

	return h.Respond(c, st, b)
}

// SearchServicesHandler : Finds all services
func SearchServicesHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/get")
	if st == 200 {
		p := h.GetSearchFilter(c)
		st, b = services.Search(au, p)
	}

	return h.Respond(c, st, b)
}

// SyncServiceHandler : Respons to POST /services/:service/sync/ and synchronizes a service with
// its provider representation
func SyncServiceHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/sync")
	if st == 200 {
		name := c.Param("name")
		st, b = services.Sync(au, name)
	}

	return h.Respond(c, st, b)
}

// ResetServiceHandler : Respons to POST /services/:service/reset/ and updates the
// service status to errored from in_progress
func ResetServiceHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/reset")
	if st == 200 {
		name := c.Param("service")
		st, b = services.Reset(au, name)
	}

	return h.Respond(c, st, b)
}

// CreateServiceHandler : Will receive a service application
func CreateServiceHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/create")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	input, definition, jsonbody, err := mapInputService(c)
	if err != nil {
		return h.Respond(c, 400, []byte(err.Error()))
	}
	isAnImport := strings.Contains(c.Path(), "/import/")
	dry := c.QueryParam("dry")
	st, b = services.Create(au, input, definition, jsonbody, isAnImport, dry)

	return h.Respond(c, st, b)
}

// UpdateServiceHandler : Not implemented
func UpdateServiceHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/update")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	name := c.Param("name")
	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = services.Update(au, name, body)
	}

	return h.Respond(c, st, b)
}

// DeleteServiceHandler : Deletes a service by name
func DeleteServiceHandler(c echo.Context) error {
	return genericDelete(c, "service", services.Delete)
}

// ForceServiceDeletionHandler : Deletes a service by name forcing it
func ForceServiceDeletionHandler(c echo.Context) error {
	return genericDelete(c, "service", services.ForceDeletion)
}
