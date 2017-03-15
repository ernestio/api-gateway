/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetLoggersHandler : responds to GET /loggers/ with a list of all
// loggers
func GetLoggersHandler(c echo.Context) (err error) {
	var loggers []models.Logger
	var body []byte
	var logger models.Logger

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	if err = logger.FindAll(&loggers); err != nil {
		return err
	}

	if body, err = json.Marshal(loggers); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// CreateLoggerHandler : responds to POST /loggers/ by creating a logger
// on the data store
func CreateLoggerHandler(c echo.Context) (err error) {
	var l models.Logger
	var body []byte

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}

	if l.Map(data) != nil {
		return h.ErrBadReqBody
	}

	if err = l.Save(); err != nil {
		return c.JSONBlob(400, []byte(err.Error()))
	}

	if body, err = json.Marshal(l); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// DeleteLoggerHandler : responds to DELETE /loggers/:id: by deleting an
// existing logger
func DeleteLoggerHandler(c echo.Context) (err error) {
	var l models.Logger

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}

	if l.Map(data) != nil {
		return h.ErrBadReqBody
	}

	if err := l.Delete(); err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}
