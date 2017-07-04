/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helpers

import (
	"io/ioutil"
	"strconv"

	"github.com/labstack/echo"
)

// GetParamFilter : Returns a filter based on parameters defined on the url stem
func GetParamFilter(c echo.Context) map[string]interface{} {
	query := make(map[string]interface{})

	fields := []string{"group", "user", "datacenter"}

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

// GetSearchFilter : Returns a filter based on url query values from the request
func GetSearchFilter(c echo.Context) map[string]interface{} {
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

// GetRequestBody : Get the request body
func GetRequestBody(c echo.Context) ([]byte, error) {
	data, err := ioutil.ReadAll(c.Request().Body)
	return data, err
}
