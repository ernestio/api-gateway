/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/loggers"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetLoggersHandler : responds to GET /loggers/ with a list of all
// loggers
func GetLoggersHandler(c echo.Context) (err error) {
	return genericList(c, "logger", loggers.List)
}

// CreateLoggerHandler : responds to POST /loggers/ by creating a logger
// on the data store
func CreateLoggerHandler(c echo.Context) (err error) {
	return genericCreate(c, "logger", loggers.Create)
}

// DeleteLoggerHandler : responds to DELETE /loggers/:id: by deleting an
// existing logger
func DeleteLoggerHandler(c echo.Context) (err error) {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = loggers.Delete(au, body)
	}

	return c.JSONBlob(s, b)
}
