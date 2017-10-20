/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/builds"
	"github.com/ernestio/api-gateway/controllers/envs"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetEnvsHandler : responds to GET /envs/envs/ with a list ahorized envs
func GetEnvsHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	project := c.Param("project")
	st, b := h.IsAuthorized(&au, "envs/get")
	if st == 200 {
		if project != "" {
			st, b = envs.List(au, &project)
		} else {
			st, b = envs.List(au, nil)
		}
	}

	return c.JSONBlob(st, b)
}

// GetEnvHandler : responds to GET /envs/:env with the
// details of an existing env
func GetEnvHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "envs/get")
	if st == 200 {
		st, b = envs.Get(au, envName(c))
	}

	return c.JSONBlob(st, b)
}

// SearchEnvsHandler : Finds all envs
func SearchEnvsHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "envs/search")
	if st == 200 {
		p := h.GetSearchFilter(c)
		st, b = envs.Search(au, p)
	}

	return h.Respond(c, st, b)
}

// CreateEnvHandler : Will receive a env application
func CreateEnvHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "envs/create")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = envs.Create(au, c.Param("project"), body)
	}

	return h.Respond(c, st, b)
}

// UpdateEnvHandler : Not implemented
func UpdateEnvHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "envs/update")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = envs.Update(au, envName(c), body)
	}

	return h.Respond(c, st, b)
}

// DeleteEnvHandler : Deletes a env by name
func DeleteEnvHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "envs/delete")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st, b = builds.Delete(au, envName(c))

	return h.Respond(c, st, b)
}

// ForceEnvDeletionHandler : Deletes an env by name forcing it
func ForceEnvDeletionHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "envs/delete")
	if st == 200 {
		st, b = envs.ForceDeletion(au, envName(c))
	}

	return h.Respond(c, st, b)
}
