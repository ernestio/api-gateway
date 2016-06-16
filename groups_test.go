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

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mockGroups = []Group{
		Group{
			ID:   1,
			Name: "test",
		},
		Group{
			ID:   1,
			Name: "test",
		},
	}
)

func getGroupSubcriber() {
	n.Subscribe("group.get", func(msg *nats.Msg) {
		var qu Group

		if len(msg.Data) > 0 {
			json.Unmarshal(msg.Data, &qu)

			for _, group := range mockGroups {
				if group.ID == qu.ID || group.Name == qu.Name {
					data, _ := json.Marshal(group)
					n.Publish(msg.Reply, data)
					return
				}
			}
		}

		n.Publish(msg.Reply, []byte(`{"error":"not found"}`))
	})
}

func findGroupSubcriber() {
	n.Subscribe("group.find", func(msg *nats.Msg) {
		var qu Group
		var ur []Group

		if len(msg.Data) == 0 {
			data, _ := json.Marshal(mockGroups)
			n.Publish(msg.Reply, data)
			return
		}

		json.Unmarshal(msg.Data, &qu)

		for _, group := range mockGroups {
			if group.Name == qu.Name || group.ID == qu.ID {
				ur = append(ur, group)
			}
		}

		data, _ := json.Marshal(ur)
		n.Publish(msg.Reply, data)
	})
}

func setGroupSubcriber() {
	n.Subscribe("group.set", func(msg *nats.Msg) {
		var u Group

		json.Unmarshal(msg.Data, &u)
		if u.ID == 0 {
			u.ID = 3
		}

		data, _ := json.Marshal(u)
		n.Publish(msg.Reply, data)
	})
}

func deleteGroupSubcriber() {
	n.Subscribe("group.del", func(msg *nats.Msg) {
		var u Datacenter

		json.Unmarshal(msg.Data, &u)

		n.Publish(msg.Reply, []byte{})
	})
}

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

					ft := jwt.New(jwt.SigningMethodHS256)
					ft.Claims["username"] = "test"
					ft.Claims["admin"] = true
					ft.Claims["group_id"] = 1

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

					ft := jwt.New(jwt.SigningMethodHS256)
					ft.Claims["username"] = "test"
					ft.Claims["admin"] = false
					ft.Claims["group_id"] = 1

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

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "test"
				ft.Claims["admin"] = true
				ft.Claims["group_id"] = 1

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
