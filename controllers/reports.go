/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/usages"
	"github.com/labstack/echo"
)

// GetUsageReportHandler : ...
func GetUsageReportHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	f := c.QueryParam("from")
	t := c.QueryParam("to")

	s, b := usages.Report(au, f, t)

	return c.JSONBlob(s, b)
}
