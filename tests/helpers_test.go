/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetParamFilter(t *testing.T) {
	e := echo.New()
	req := new(http.Request)
	rec := httptest.NewRecorder()
	Convey("Scenario: getting an empty http context", t, func() {
		c := e.NewContext(req, rec)
		Convey("when it is converted to a query", func() {
			query := h.GetParamFilter(c)
			Convey("the query is also empty", func() {
				So(query, ShouldNotBeNil)
				So(len(query), ShouldEqual, 0)
			})
		})
	})

	Convey("Scenario: getting an http context", t, func() {
		c := e.NewContext(req, rec)
		Convey("when it has a name paramater", func() {
			c.SetParamNames("name")
			c.SetParamValues("somename")
			Convey("and it is converted to a query", func() {
				query := h.GetParamFilter(c)
				Convey("the query has the name and its value", func() {
					So(query, ShouldNotBeNil)
					So(len(query), ShouldEqual, 1)
					So(query["name"], ShouldEqual, "somename")
				})
			})
		})
	})

	Convey("Scenario: getting an http context", t, func() {
		c := e.NewContext(req, rec)
		Convey("when it has a service paramater", func() {
			c.SetParamNames("service")
			c.SetParamValues("somename")
			Convey("and it is converted to a query", func() {
				query := h.GetParamFilter(c)
				Convey("the query has the name and its value", func() {
					So(query, ShouldNotBeNil)
					So(len(query), ShouldEqual, 1)
					So(query["name"], ShouldEqual, "somename")
				})
			})
		})
	})

	Convey("Scenario: getting an http context", t, func() {
		c := e.NewContext(req, rec)
		Convey("when it has a build paramater", func() {
			c.SetParamNames("build")
			c.SetParamValues("someid")
			Convey("and it is converted to a query", func() {
				query := h.GetParamFilter(c)
				Convey("the query has the id and its value", func() {
					So(query, ShouldNotBeNil)
					So(len(query), ShouldEqual, 1)
					So(query["id"], ShouldEqual, "someid")
				})
			})
		})
	})

}
