/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func getUsageReportHandler(c echo.Context) (err error) {
	var usage Usage
	var reportables []Usage
	var body []byte
	var from, to int64

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
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

	if body, err = renderUsageReport(reportables); err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, body)
}
