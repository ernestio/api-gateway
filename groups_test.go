/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGroups(t *testing.T) {
	Convey("Given group handler", t, func() {
		// setup nats connection
		os.Setenv("JWT_SECRET", "test")
		setup()

		Convey("When getting a list of groups", func() {
			findGroupSubcriber()

			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
			c.SetPath("/groups/")

			Convey("It should return the correct set of data", func() {
				var u []Group

				err := getGroupsHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &u)

				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, 200)
				So(len(u), ShouldEqual, 2)
				So(u[0].ID, ShouldEqual, 1)
				So(u[0].Name, ShouldEqual, "test")
			})

		})

		Convey("When getting a single group", func() {
			getGroupSubcriber()

			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			c.SetPath("/groups/:group")
			c.SetParamNames("group")
			c.SetParamValues("test")

			Convey("It should return the correct set of data", func() {
				var u Group

				err := getGroupHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &u)

				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, 200)
				So(u.ID, ShouldEqual, 1)
				So(u.Name, ShouldEqual, "test")
			})
		})

		Convey("When creating a group", func() {
			setGroupSubcriber()

			Convey("With a valid payload", func() {
				data, _ := json.Marshal(Group{Name: "new-test"})

				Convey("As an admin user", func() {
					e := echo.New()
					req, _ := http.NewRequest("POST", "/groups/", bytes.NewReader(data))
					rec := httptest.NewRecorder()
					c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

					ft := generateTestToken(1, "test", true)

					c.SetPath("/groups/")
					c.Set("user", ft)

					Convey("It should create the group and return the correct set of data", func() {
						var u Group

						err := createGroupHandler(c)
						So(err, ShouldBeNil)

						resp := rec.Body.Bytes()
						err = json.Unmarshal(resp, &u)

						So(err, ShouldBeNil)
						So(u.ID, ShouldEqual, 3)
						So(u.Name, ShouldEqual, "new-test")
					})
				})

				Convey("As an non-admin user", func() {
					e := echo.New()
					req, _ := http.NewRequest("POST", "/groups/", bytes.NewReader(data))
					rec := httptest.NewRecorder()
					c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

					ft := generateTestToken(1, "test", false)

					c.SetPath("/groups/")
					c.Set("user", ft)

					Convey("It should return with 403 unauthorized", func() {
						err := createGroupHandler(c)
						So(err, ShouldNotBeNil)
						So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
					})
				})
			})

			Convey("With an invalid payload", func() {
				data := []byte(`{"incorrect_name": "fail"}`)

				e := echo.New()
				req, _ := http.NewRequest("POST", "/groups/", bytes.NewReader(data))
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := generateTestToken(1, "test", true)

				c.Set("user", ft)
				c.SetPath("/groups/")

				Convey("It should error with 400 bad request", func() {
					err := createGroupHandler(c)
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
				})
			})
		})

	})
}
