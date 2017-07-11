/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/services"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetServiceBuildHandler : gets the details of a specific service build
func GetServiceBuildHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	query := h.GetAuthorizedParamFilter(c, &au)
	s, b := services.GetBuild(au, query)

	return c.JSONBlob(s, b)
}

// DelServiceBuildHandler : will delete the specified build from a service
func DelServiceBuildHandler(c echo.Context) (err error) {
	return genericDelete(c, "build", services.DelBuild)
}
