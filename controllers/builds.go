/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/envs"
	"github.com/labstack/echo"
)

// DelServiceBuildHandler : will delete the specified build from a service
func DelServiceBuildHandler(c echo.Context) (err error) {
	return genericDelete(c, "build", envs.DelBuild)
}
