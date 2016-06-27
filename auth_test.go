/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	. "github.com/smartystreets/goconvey/convey"
)

var mockToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDU4ODUwMTE5MSwiZ3JvdXBfaWQiOjIsInVzZXJuYW1lIjoidGVzdDIifQ.SrP29afiIPjtIbdKrUXyf9B8m6_fPVTI0mgH6s4Y_VY"

func TestAuth(t *testing.T) {
	Convey("Given the auth handler", t, func() {
		testsSetup()
		setup()

		Convey("When attempting to login", func() {
			getUserSubscriber(1)

			Convey("With valid credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				req.PostForm = url.Values{"username": {"test2"}, "password": {"test2"}}

				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
				c.SetPath("/auth/")

				err := authenticate(c)
				resp := rec.Body.String()

				Convey("It should return a jwt token", func() {
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusOK)
					So(strings.Contains(resp, "token"), ShouldBeTrue)
				})
			})

			Convey("With invalid credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				req.PostForm = url.Values{"username": {"test2"}, "password": {"wrong"}}
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
				c.SetPath("/auth/")

				err := authenticate(c)
				resp := rec.Body.String()

				Convey("It should not return a jwt token and error", func() {
					So(err, ShouldNotBeNil)
					So(strings.Contains(resp, "token"), ShouldBeFalse)
				})
			})

			Convey("With no credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
				c.SetPath("/auth/")

				err := authenticate(c)
				resp := rec.Body.String()

				Convey("It should not return a jwt token and error", func() {
					So(err, ShouldNotBeNil)
					So(strings.Contains(resp, "token"), ShouldBeFalse)
				})
			})
		})
	})

	Convey("Given a protected route", t, func() {
		testsSetup()
		setup()

		Convey("When attempting to retrieve data", func() {
			getUserSubscriber(1)
			findUserSubscriber()

			Convey("With valid credentials", func() {
				e := echo.New()
				req := new(http.Request)
				authHeader := fmt.Sprintf("Bearer %s", mockToken)
				req.Header = http.Header{}
				req.Header.Add("Authorization", authHeader)
				rec := httptest.NewRecorder()

				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
				c.SetPath("/users/")
				h := middleware.JWT([]byte(secret))(getUsersHandler)

				Convey("It should return the correct data", func() {
					err := h(c)
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusOK)
					So(strings.Contains(rec.Body.String(), "name"), ShouldBeTrue)
				})
			})

			Convey("With invalid credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
				c.SetPath("/users/")

				h := middleware.JWT([]byte(secret))(getUsersHandler)

				err := h(c)
				resp := rec.Body.String()

				Convey("It should return an 400 bad request", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
					So(strings.Contains(resp, "id"), ShouldBeFalse)
				})
			})

			Convey("With no credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
				c.SetPath("/users/")
				h := middleware.JWT([]byte(secret))(getUsersHandler)

				err := h(c)
				resp := rec.Body.String()

				Convey("It should return an 400 bad request", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
					So(strings.Contains(resp, "id"), ShouldBeFalse)
				})
			})
		})
	})
}
