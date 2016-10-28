/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

// getLoggersHandler : responds to GET /loggers/ with a list of all
// loggers
func getLoggersHandler(c echo.Context) (err error) {
	var loggers []Logger
	var body []byte
	var logger Logger

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if err = logger.FindAll(&loggers); err != nil {
		return err
	}

	if body, err = json.Marshal(loggers); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// createLoggerHandler : responds to POST /loggers/ by creating a logger
// on the data store
func createLoggerHandler(c echo.Context) (err error) {
	var l Logger
	var body []byte

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if l.Map(c) != nil {
		return ErrBadReqBody
	}

	if err = l.Save(); err != nil {
		return c.JSONBlob(400, []byte(err.Error()))
	}

	if body, err = json.Marshal(l); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// deleteLoggerHandler : responds to DELETE /loggers/:id: by deleting an
// existing logger
func deleteLoggerHandler(c echo.Context) (err error) {
	var l Logger

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if l.Map(c) != nil {
		return ErrBadReqBody
	}

	if err := l.Delete(); err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}
