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

// GetServicesHandler : responds to GET /services/ with a list of all
// services for current user group
func GetServicesHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	s, b := services.List(au)

	return c.JSONBlob(s, b)
}

// GetServiceBuildsHandler : gets the list of builds for the specified
// service
func GetServiceBuildsHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	p := h.GetParamFilter(c)
	s, b := services.Builds(au, p)

	return c.JSONBlob(s, b)
}

// GetServiceHandler : responds to GET /services/:service with the
// details of an existing service
func GetServiceHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	p := h.GetParamFilter(c)
	s, b := services.Get(au, p)

	return c.JSONBlob(s, b)
}

// SearchServicesHandler : Finds all services
func SearchServicesHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	p := h.GetSearchFilter(c)
	s, b := services.Search(au, p)

	return c.JSONBlob(s, b)
}

// SyncServiceHandler : Respons to POST /services/:service/sync/ and synchronizes a service with
// its provider representation
func SyncServiceHandler(c echo.Context) error {
	if err := h.Licensed(); err != nil {
		return err
	}
	name := c.Param("name")
	au := AuthenticatedUser(c)
	s, b := services.Sync(au, name)

	return c.JSONBlob(s, b)
}

// ResetServiceHandler : Respons to POST /services/:service/reset/ and updates the
// service status to errored from in_progress
func ResetServiceHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	name := c.Param("service")
	s, b := services.Reset(au, name)

	return c.JSONBlob(s, b)
}

// CreateUUIDHandler : Creates an unique id
func CreateUUIDHandler(c echo.Context) error {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = services.UUID(au, body)
	}

	return c.JSONBlob(s, b)
}

// CreateServiceHandler : Will receive a service application
func CreateServiceHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	input, definition, jsonbody, err := mapInputService(c)
	if err != nil {
		return c.JSONBlob(400, []byte(err.Error()))
	}
	isAnImport := strings.Contains(c.Path(), "/import/")
	dry := c.QueryParam("dry")
	s, b := services.CreateServiceHandler(au, input, definition, jsonbody, isAnImport, dry)

	return c.JSONBlob(s, b)
}

// UpdateServiceHandler : Not implemented
func UpdateServiceHandler(c echo.Context) error {
	if err := h.Licensed(); err != nil {
		return err
	}
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)

	name := c.Param("name")
	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = services.Update(au, name, body)
	}

	return c.JSONBlob(s, b)
}

// DeleteServiceHandler : Deletes a service by name
func DeleteServiceHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	name := c.Param("name")
	s, b := services.Delete(au, name)

	return c.JSONBlob(s, b)
}

// ForceServiceDeletionHandler : Deletes a service by name forcing it
func ForceServiceDeletionHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	name := c.Param("name")
	s, b := services.ForceDeletion(au, name)

	return c.JSONBlob(s, b)
}
