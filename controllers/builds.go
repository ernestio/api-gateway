/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/builds"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetBuildHandler : gets the details of a specific env build
func GetBuildHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/build")
	if st == 200 {
		st, b = builds.Get(au, c.Param("build"))
	}

	return c.JSONBlob(st, b)
}

// GetBuildsHandler : gets the list of builds for the specified
// env
func GetBuildsHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/builds")
	if st == 200 {
		st, b = builds.List(au, buildID(c))
	}

	return h.Respond(c, st, b)
}

// GetBuildMappingHandler : gets the mapping of a build
func GetBuildMappingHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/builds")
	if st == 200 {
		st, b = builds.Mapping(au, buildID(c))
	}

	return h.Respond(c, st, b)
}

// GetBuildDefinitionHandler : gets the mapping of a build
func GetBuildDefinitionHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/builds")
	if st == 200 {
		st, b = builds.Mapping(au, buildID(c))
	}

	return h.Respond(c, st, b)
}

// CreateBuildHandler : Will receive a env application
func CreateBuildHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/create")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	definition, raw, err := mapInputBuild(c)
	if err != nil {
		return h.Respond(c, 400, []byte(err.Error()))
	}

	dry := c.QueryParam("dry")
	st, b = builds.Create(au, &definition, raw, dry)

	return h.Respond(c, st, b)
}

/*

// DelEnvBuildHandler : will delete the specified build from a service
func DelEnvBuildHandler(c echo.Context) (err error) {
	return genericDelete(c, "build", envs.DelBuild)
}
*/

func buildID(c echo.Context) string {
	env := c.Param("env")
	proj := c.Param("project")
	return proj + models.EnvNameSeparator + env
}
