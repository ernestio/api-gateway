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

	"github.com/ernestio/api-gateway/config"
	"github.com/ernestio/api-gateway/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	. "github.com/smartystreets/goconvey/convey"
)

var mockToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDU4ODUwMTE5MSwiZ3JvdXBfaWQiOjIsInVzZXJuYW1lIjoidGVzdDIifQ.SrP29afiIPjtIbdKrUXyf9B8m6_fPVTI0mgH6s4Y_VY"

func TestAuth(t *testing.T) {
	Convey("Given the auth handler", t, func() {
		testsSetup()
		config.Setup()

		Convey("When attempting to login", func() {
			getUserSubscriber(1)
			getGroupSubscriber(1)

			Convey("With valid credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				req.PostForm = url.Values{"username": {"test2"}, "password": {"test1234"}}

				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/auth/")

				err := controllers.AuthenticateHandler(c)
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

				req.PostForm = url.Values{"username": {"test2"}, "password": {"wrong1234"}}
				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/auth/")

				err := controllers.AuthenticateHandler(c)
				resp := rec.Body.String()

				Convey("It should not return a jwt token and error", func() {
					So(err, ShouldNotBeNil)
					So(strings.Contains(resp, "token"), ShouldBeFalse)
				})
			})

			Convey("With a username using invalid characters", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				req.PostForm = url.Values{"username": {"test^2"}, "password": {"test1234"}}
				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/auth/")

				err := controllers.AuthenticateHandler(c)
				resp := rec.Body.String()

				Convey("It should not return a jwt token and error", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
					So(err.(*echo.HTTPError).Message, ShouldEqual, "Username can only contain the following characters: a-z 0-9 @._-")
					So(resp, ShouldNotContainSubstring, "token")
				})
			})

			Convey("With a password using invalid characters", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				req.PostForm = url.Values{"username": {"test2"}, "password": {"test^1234"}}
				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/auth/")

				err := controllers.AuthenticateHandler(c)
				resp := rec.Body.String()

				Convey("It should not return a jwt token and error", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
					So(err.(*echo.HTTPError).Message, ShouldEqual, "Password can only contain the following characters: a-z 0-9 @._-")
					So(resp, ShouldNotContainSubstring, "token")
				})
			})

			Convey("With no username", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				req.PostForm = url.Values{"username": {""}, "password": {"test"}}
				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/auth/")

				err := controllers.AuthenticateHandler(c)
				resp := rec.Body.String()

				Convey("It should not return a jwt token and error", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
					So(err.(*echo.HTTPError).Message, ShouldEqual, "Username cannot be empty")
					So(resp, ShouldNotContainSubstring, "token")
				})
			})

			Convey("With no password", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				req.PostForm = url.Values{"username": {"test2"}, "password": {""}}
				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/auth/")

				err := controllers.AuthenticateHandler(c)
				resp := rec.Body.String()

				Convey("It should not return a jwt token and error", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
					So(err.(*echo.HTTPError).Message, ShouldEqual, "Password cannot be empty")
					So(resp, ShouldNotContainSubstring, "token")
				})
			})

			Convey("With no credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/auth/")

				err := controllers.AuthenticateHandler(c)
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
		config.Setup()

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

				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/users/")
				h := middleware.JWT([]byte(controllers.Secret))(controllers.GetUsersHandler)

				Convey("It should return the correct data", func() {
					err := h(c)
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusOK)
					So(rec.Body.String(), ShouldContainSubstring, "name")
				})
			})

			Convey("With invalid credentials", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()

				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/users/")

				h := middleware.JWT([]byte(controllers.Secret))(controllers.GetUsersHandler)

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

				c := e.NewContext(req, echo.NewResponse(rec, e))
				c.SetPath("/users/")
				h := middleware.JWT([]byte(controllers.Secret))(controllers.GetUsersHandler)

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
