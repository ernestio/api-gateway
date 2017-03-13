/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"log"
	"net/http"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
	"github.com/labstack/echo"
)

// GetUsageReportHandler : ...
func GetUsageReportHandler(c echo.Context) (err error) {
	var usage models.Usage
	var reportables []models.Usage
	var body []byte
	var from, to int64

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	layout := "2006-01-02"

	if c.QueryParam("from") != "" {
		fromTime, err := time.Parse(layout, c.QueryParam("from"))
		if err != nil {
			log.Println("Invalid from date on usage report")
			log.Println(err.Error())
			return err
		}
		from = fromTime.Unix()
	}
	if c.QueryParam("to") != "" {
		toTime, err := time.Parse(layout, c.QueryParam("to"))
		if err != nil {
			log.Println("Invalid to date on usage report")
			log.Println(err.Error())
			return err
		}
		to = toTime.Unix()
	}

	if err = usage.FindAllInRange(from, to, &reportables); err != nil {
		return err
	}

	if body, err = views.RenderUsageReport(reportables); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, body)
}
