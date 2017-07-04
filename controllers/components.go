/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"strings"

	"github.com/ernestio/api-gateway/controllers/components"
	"github.com/labstack/echo"
)

// GetAllComponentsHandler : ...
func GetAllComponentsHandler(c echo.Context) (err error) {
	parts := strings.Split(c.Path(), "/")
	component := parts[len(parts)-2] + "s"
	d := c.QueryParam("datacenter")
	s := c.QueryParam("service")

	st, b := components.List(d, s, component)

	return c.JSONBlob(st, b)
}
