/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	. "github.com/smartystreets/goconvey/convey"
)

var mockToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDU4NzgyOTc4MSwibmFtZSI6InRlc3QyIiwidXNlcm5hbWUiOiJ0ZXN0MiJ9.mRYxkFlLrjJkV49EAQ4wnkk4diNUwl3yJxWpZLXEHFE"

func TestAuth(t *testing.T) {
	Convey("Given the auth handler", t, func() {
		// setup nats connection
		os.Setenv("JWT_SECRET", "test")
		setup()

		Convey("When attempting to login", func() {
			findUserSubcriber()

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
		// setup nats connection
		os.Setenv("JWT_SECRET", "test")
		setup()

		Convey("When attempting to retrieve data", func() {
			findUserSubcriber()
			getUsersSubcriber()

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
				err := h(c)
				resp := rec.Body.String()

				Convey("It should return the correct data", func() {
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusOK)
					So(strings.Contains(resp, "name"), ShouldBeTrue)
				})
			})

			Convey("With invalid credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
				c.SetPath("/users/")
				h := middleware.JWT([]byte(secret))(getUsersHandler)
				err := h(c).(*echo.HTTPError)
				resp := rec.Body.String()

				Convey("It should return an unauthorized error", func() {
					So(err, ShouldNotBeNil)
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
				err := h(c).(*echo.HTTPError)
				resp := rec.Body.String()

				Convey("It should return an unauthorized error", func() {
					So(err, ShouldNotBeNil)
					So(strings.Contains(resp, "id"), ShouldBeFalse)
				})
			})
		})
	})
}
