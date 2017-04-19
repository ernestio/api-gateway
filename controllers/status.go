/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

// GetStatusHandler : responds to GET /status/
func GetStatusHandler(c echo.Context) (err error) {
	return c.JSONBlob(http.StatusOK, []byte(`"success"`))
}
