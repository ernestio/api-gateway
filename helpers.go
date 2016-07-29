/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var (
	// ErrUnauthorized : HTTP 403 error
	ErrUnauthorized = echo.NewHTTPError(http.StatusForbidden, "")
	// ErrNotFound : HTTP 404 error
	ErrNotFound = echo.NewHTTPError(http.StatusNotFound, "")
	// ErrBadReqBody : HTTP 400 error
	ErrBadReqBody = echo.NewHTTPError(http.StatusBadRequest, "")
	// ErrGatewayTimeout : HTTP 504 error
	ErrGatewayTimeout = echo.NewHTTPError(http.StatusGatewayTimeout, "")
	// ErrInternal : HTTP 500 error
	ErrInternal = echo.NewHTTPError(http.StatusInternalServerError, "")
	// ErrNotImplemented : HTTP 405 error
	ErrNotImplemented = echo.NewHTTPError(http.StatusNotImplemented, "")
	// ErrExists : HTTP Error
	ErrExists = echo.NewHTTPError(http.StatusSeeOther, "")
)

// Get the authenticated user from the JWT Token
func authenticatedUser(c echo.Context) User {
	var u User

	user := c.Get("user").(*jwt.Token)

	claims, ok := user.Claims.(jwt.MapClaims)
	if ok {
		u.Username = claims["username"].(string)
		u.GroupID = int(claims["group_id"].(float64))
		u.Admin = claims["admin"].(bool)
	}

	return u
}

// Returns a filter based on parameters defined on the url stem
func getParamFilter(c echo.Context) map[string]interface{} {
	query := make(map[string]interface{})

	fields := []string{"group", "user", "group", "datacenter"}

	// Process ID's as int's
	for _, field := range fields {
		if val := c.Param(field); val != "" {
			id, err := strconv.Atoi(val)
			if err == nil {
				query["id"] = id
			}
		}
	}

	if c.Param("name") != "" {
		query["name"] = c.Param("name")
	}

	if c.Param("service") != "" {
		query["name"] = c.Param("service")
	}

	if c.Param("build") != "" {
		query["id"] = c.Param("build")
	}

	return query
}

// Returns a filter based on url query values from the request
func getSearchFilter(c echo.Context) map[string]interface{} {
	query := make(map[string]interface{})

	fields := []string{"id", "user_id", "group_id", "datacenter_id", "service_id"}

	// Process ID's as int's
	for _, field := range fields {
		if val := c.QueryParam(field); val != "" {
			id, err := strconv.Atoi(val)
			if err == nil {
				query[field] = id
			}
		}
	}

	if c.QueryParam("name") != "" {
		query["name"] = c.QueryParam("name")
	}

	return query
}
