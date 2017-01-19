/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	// "log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// getAllComponentsHandler : ...
func getAllComponentsHandler(c echo.Context) (err error) {
	var body []byte
	var d Datacenter

	parts := strings.Split(c.Path(), "/")
	component := parts[len(parts)-2] + "s"

	if err := d.FindByName(c.QueryParam("datacenter"), &d); err != nil {
		return err
	}

	query := make(map[string]interface{})
	query["expects_response"] = true
	query["aws_access_key_id"] = d.AccessKeyID
	query["aws_secret_access_key"] = d.SecretAccessKey
	query["datacenter_region"] = d.Region
	service := c.QueryParam("service")
	if service != "" {
		tags := make(map[string]string)
		tags["ernest.service"] = c.QueryParam("service")
		query["tags"] = tags
	}

	components := make(map[string]interface{})
	if err = NewBaseModel(component).callStoreBy("find.aws", query, &components); err != nil {
		return c.JSONBlob(500, []byte("An internal error occured"))
	}

	if components["components"] == nil {
		return c.JSONBlob(200, []byte("[]"))
	}

	list := components["components"].([]interface{})
	for i := range list {
		component := list[i].(map[string]interface{})
		delete(component, "_uuid")
		delete(component, "_batch_id")
		delete(component, "_type")
		delete(component, "aws_access_key_id")
		delete(component, "aws_secret_access_key")
		delete(component, "_uuid")
	}

	if body, err = json.Marshal(list); err != nil {
		return c.JSONBlob(500, []byte("Oops, somethign went wrong"))
	}

	return c.JSONBlob(http.StatusOK, body)
}
