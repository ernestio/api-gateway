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
	os.Setenv("JWT_SECRET", "test")
	setup()

	Convey("Scenario: getting a list of groups", t, func() {
		Convey("Given groups exist on the store", func() {
			findGroupSubcriber()
			Convey("When I call /groups/", func() {
				resp, err := doRequest("GET", "/groups/", nil, nil, getGroupsHandler, nil)
				Convey("Then I should have a response existing groups", func() {
					var d []Group
					So(err, ShouldBeNil)

					err = json.Unmarshal(resp, &d)

					So(err, ShouldBeNil)
					So(len(d), ShouldEqual, 2)
					So(d[0].ID, ShouldEqual, 1)
					So(d[0].Name, ShouldEqual, "test")
				})
			})

			SkipConvey("Given no groups on the store", func() {
			})
		})
	})

	Convey("Scenario: getting a single group", t, func() {
		Convey("Given the group exist on the store", func() {
			getGroupSubcriber()
			Convey("And I call /groups/:group on the api", func() {
				params := make(map[string]string)
				params["group"] = "1"
				resp, err := doRequest("GET", "/groups/:group", params, nil, getGroupHandler, nil)

				Convey("When I'm authenticated as admin user", func() {
					Convey("Then I should get the existing group", func() {
						var g Group

						So(err, ShouldBeNil)

						err = json.Unmarshal(resp, &g)

						So(err, ShouldBeNil)
						So(g.ID, ShouldEqual, 1)
						So(g.Name, ShouldEqual, "test")
					})
				})
			})
		})
	})

	Convey("Given group handler", t, func() {

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

	Convey("Scenario: deleting a group", t, func() {
		Convey("Given a group exists on the store", func() {
			deleteGroupSubcriber()

			Convey("When I call DELETE /groups/:group", func() {
				ft := generateTestToken(1, "test", false)

				params := make(map[string]string)
				params["group"] = "test"
				_, err := doRequest("DELETE", "/groups/:group", params, nil, deleteGroupHandler, ft)

				Convey("It should delete the group and return ok", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
}
