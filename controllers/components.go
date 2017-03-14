/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetAllComponentsHandler : ...
func GetAllComponentsHandler(c echo.Context) (err error) {
	var body []byte
	var d models.Datacenter

	parts := strings.Split(c.Path(), "/")
	component := parts[len(parts)-2] + "s"

	if err := d.FindByName(c.QueryParam("datacenter"), &d); err != nil {
		return err
	}

	tags := make(map[string]string)
	service := c.QueryParam("service")
	if service != "" {
		tags["ernest.service"] = c.QueryParam("service")
	}
	aws := models.AWSComponent{
		Datacenter: &d,
		Name:       component,
		Tags:       tags,
	}
	list, err := aws.FindBy()
	if err != nil {
		return c.JSONBlob(500, []byte("An internal error occured"))
	}

	if body, err = json.Marshal(list); err != nil {
		return c.JSONBlob(500, []byte("Oops, somethign went wrong"))
	}

	return c.JSONBlob(http.StatusOK, body)
}
