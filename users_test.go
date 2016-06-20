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

func TestUsers(t *testing.T) {
	Convey("Scenario: getting a list of users", t, func() {
		// setup nats connection
		os.Setenv("JWT_SECRET", "test")
		setup()

		findUserSubcriber()

		e := echo.New()
		req := new(http.Request)
		rec := httptest.NewRecorder()
		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		c.SetPath("/users/")

		Convey("It should return the correct set of data", func() {
			var u []User

			err := getUsersHandler(c)
			So(err, ShouldBeNil)

			resp := rec.Body.Bytes()
			err = json.Unmarshal(resp, &u)

			So(err, ShouldBeNil)
			So(rec.Code, ShouldEqual, 200)
			So(len(u), ShouldEqual, 2)
			So(u[0].ID, ShouldEqual, 1)
			So(u[0].Username, ShouldEqual, "test")
		})

	})

	Convey("Scenario: getting a single user", t, func() {
		getUserSubcriber()

		e := echo.New()
		req := new(http.Request)
		rec := httptest.NewRecorder()
		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		c.SetPath("/users/:user")
		c.SetParamNames("user")
		c.SetParamValues("test")

		Convey("It should return the correct set of data", func() {
			var u User

			err := getUserHandler(c)
			So(err, ShouldBeNil)

			resp := rec.Body.Bytes()
			err = json.Unmarshal(resp, &u)

			So(err, ShouldBeNil)
			So(rec.Code, ShouldEqual, 200)
			So(u.ID, ShouldEqual, 1)
			So(u.Username, ShouldEqual, "test")
		})
	})

	Convey("Scenario: creating a user", t, func() {
		setUserSubcriber()

		data, _ := json.Marshal(User{GroupID: 1, Username: "new-test", Password: "test"})

		Convey("As an admin user", func() {
			e := echo.New()
			req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(data))
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			ft := generateTestToken(1, "test", true)

			c.SetPath("/users/")
			c.Set("user", ft)

			Convey("It should create the user and return the correct set of data", func() {
				var u User

				err := createUserHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &u)

				So(err, ShouldBeNil)
				So(u.ID, ShouldEqual, 3)
				So(u.Username, ShouldEqual, "new-test")
			})
		})

		Convey("As an non-admin user", func() {
			e := echo.New()
			req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(data))
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			ft := generateTestToken(1, "test", false)

			c.SetPath("/users/")
			c.Set("user", ft)

			Convey("It should return with 403 unauthorized", func() {
				err := createUserHandler(c)
				So(err, ShouldNotBeNil)
				So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
			})
		})

		Convey("With an invalid payload", func() {
			data := []byte(`{"group_id": 1, "username": "fail"}`)

			e := echo.New()
			req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(data))
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			ft := generateTestToken(1, "test", true)

			c.Set("user", ft)
			c.SetPath("/users/")

			Convey("It should error with 400 bad request", func() {
				err := createUserHandler(c)
				So(err, ShouldNotBeNil)
				So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
			})

		})

	})

	Convey("Scenario: updating a user", t, func() {
		setUserSubcriber()

		Convey("As an admin user", func() {
			data, _ := json.Marshal(User{GroupID: 1, ID: 1, Username: "test2", Password: "test"})

			e := echo.New()
			req, _ := http.NewRequest("POST", "/users/test", bytes.NewReader(data))
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			ft := generateTestToken(1, "admin", true)

			c.SetPath("/users/:user")
			c.SetParamNames("user")
			c.SetParamValues("test")
			c.Set("user", ft)

			Convey("It should update the user and return the correct set of data", func() {
				var u User

				err := updateUserHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &u)

				So(err, ShouldBeNil)
				So(u.ID, ShouldEqual, 1)
				So(u.GroupID, ShouldEqual, 1)
				So(u.Username, ShouldEqual, "test2")
			})
		})

		Convey("As an non-admin user", func() {
			Convey("Where a user updates itself", func() {
				data, _ := json.Marshal(User{GroupID: 1, ID: 1, Username: "test", Password: "test2"})

				e := echo.New()
				req, _ := http.NewRequest("POST", "/users/test", bytes.NewReader(data))
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := generateTestToken(1, "test", false)

				c.SetPath("/users/:user")
				c.SetParamNames("user")
				c.SetParamValues("test")
				c.Set("user", ft)

				Convey("It should update the user and return the correct set of data", func() {
					var u User

					err := updateUserHandler(c)
					So(err, ShouldBeNil)

					resp := rec.Body.Bytes()
					err = json.Unmarshal(resp, &u)

					So(err, ShouldBeNil)
					So(u.ID, ShouldEqual, 1)
					So(u.GroupID, ShouldEqual, 1)
					So(u.Username, ShouldEqual, "test")
				})
			})

			Convey("Where a non-admin user updates another user", func() {
				data, _ := json.Marshal(User{GroupID: 1, Username: "test2", Password: "test2"})

				e := echo.New()
				req, _ := http.NewRequest("POST", "/users/test2", bytes.NewReader(data))
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := generateTestToken(1, "test", false)

				c.SetPath("/users/:user")
				c.SetParamNames("user")
				c.SetParamValues("test2")
				c.Set("user", ft)

				Convey("It should return with 403 unauthorized", func() {
					err := updateUserHandler(c)
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
				})
			})
		})

		Convey("When updating a user that doesn't exist", func() {
			data := []byte(`{"group_id": 1, "username": "fake-user", "password": "fake"}`)

			e := echo.New()
			req, _ := http.NewRequest("POST", "/users/fake-user", bytes.NewReader(data))
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			ft := generateTestToken(1, "test", true)

			c.Set("user", ft)
			c.SetPath("/users/")
			c.SetParamNames("user")
			c.SetParamValues("fake-user")

			Convey("It should error with 404 doesn't exist", func() {
				err := updateUserHandler(c)
				So(err, ShouldNotBeNil)
				So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
			})
		})
	})

	Convey("Scenario: deleting a user", t, func() {
		deleteUserSubcriber()

		e := echo.New()
		req := http.Request{Method: "DELETE"}
		rec := httptest.NewRecorder()
		c := e.NewContext(standard.NewRequest(&req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		c.SetPath("/users/:user")
		c.SetParamNames("user")
		c.SetParamValues("test")

		Convey("It should delete the user and return a 200 ok", func() {
			err := deleteUserHandler(c)
			So(err, ShouldBeNil)
			So(rec.Code, ShouldEqual, 200)
		})
	})
}
