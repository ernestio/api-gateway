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

// GetEnvsHandler : responds to GET /envs/envs/ with a list ahorized services
func GetEnvsHandler(c echo.Context) (err error) {
	return genericList(c, "service", envs.List)
}

// GetEnvHandler : responds to GET /envs/:env with the
// details of an existing env
func GetEnvHandler(c echo.Context) (err error) {
	return genericGet(c, "service", envs.Get)
}

// SearchEnvsHandler : Finds all envs
func SearchEnvsHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/search")
	if st == 200 {
		p := h.GetSearchFilter(c)
		st, b = envs.Search(au, p)
	}

	return h.Respond(c, st, b)
}

// ResetEnvHandler : Respons to POST /envs/:env/reset/ and updates the
// env status to errored from in_progress
func ResetEnvHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/reset")
	if st == 200 {
		st, b = envs.Reset(au, buildID(c))
	}

	return h.Respond(c, st, b)
}

// CreateEnvHandler : Will receive a env application
func CreateEnvHandler(c echo.Context) error {
	return genericCreate(c, "service", envs.Create)
}

// UpdateEnvHandler : Not implemented
func UpdateEnvHandler(c echo.Context) error {
	return genericUpdate(c, "service", envs.Update)
}

// DeleteEnvHandler : Deletes a env by name
func DeleteEnvHandler(c echo.Context) error {
	return genericDelete(c, "service", builds.Delete)
}

// ForceEnvDeletionHandler : Deletes an env by name forcing it
func ForceEnvDeletionHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/delete")
	if st == 200 {
		st, b = envs.ForceDeletion(au, buildID(c))
	}

	return h.Respond(c, st, b)
}
