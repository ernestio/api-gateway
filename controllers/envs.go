/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/envs"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetEnvsHandler : responds to GET /envs/envs/ with a list ahorized services
func GetEnvsHandler(c echo.Context) (err error) {
	return genericList(c, "service", envs.List)
}

// GetEnvBuildHandler : gets the details of a specific env build
func GetEnvBuildHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/build")
	if st == 200 {
		query := h.GetAuthorizedParamFilter(c, &au)
		st, b = envs.GetBuild(au, query)
	}

	return c.JSONBlob(st, b)
}

// GetEnvBuildsHandler : gets the list of builds for the specified
// env
func GetEnvBuildsHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/builds")
	if st == 200 {
		st, b = envs.Builds(au, buildID(c))
	}

	return h.Respond(c, st, b)
}

// GetEnvHandler : responds to GET /envs/:env with the
// details of an existing env
func GetEnvHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/get")
	if st == 200 {
		st, b = envs.Get(au, buildID(c))
	}

	return h.Respond(c, st, b)
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

// SyncEnvHandler : Respons to POST /envs/:env/sync/ and synchronizes a env with
// its provider representation
func SyncEnvHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/sync")
	if st == 200 {
		st, b = envs.Sync(au, buildID(c))
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
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/create")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	input, definition, jsonbody, err := mapInputService(c)
	if err != nil {
		return h.Respond(c, 400, []byte(err.Error()))
	}

	input.Name = buildStringID(input.Datacenter, input.Name)

	dry := c.QueryParam("dry")
	st, b = envs.Create(au, input, definition, jsonbody, isAnImport, dry)

	return h.Respond(c, st, b)
}

// UpdateEnvHandler : Not implemented
func UpdateEnvHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/update")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = envs.Update(au, buildID(c), body)
	}

	return h.Respond(c, st, b)
}

// DeleteEnvHandler : Deletes a env by name
func DeleteEnvHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/delete")
	if st == 200 {
		st, b = envs.Delete(au, buildID(c))
	}

	return h.Respond(c, st, b)
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

// CreateEnvBuildHandler : Creates a build on an environment
func CreateEnvBuildHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/create")
	if st == 200 {
		st, b := envs.Apply()
	}

	return h.Respond(c, st, b)
}

// ImportEnvHandler : Creates an import build on an environment
func ImportEnvHandler(c echo.Context) error {

}

// DelEnvBuildHandler : will delete the specified build from a service
func DelEnvBuildHandler(c echo.Context) (err error) {
	return genericDelete(c, "build", envs.DelBuild)
}

func buildID(c echo.Context) string {
	env := c.Param("env")
	proj := c.Param("project")
	return buildStringID(proj, env)
}

func buildStringID(project, env string) string {
	return project + models.EnvNameSeparator + env
}
