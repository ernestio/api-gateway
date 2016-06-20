/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
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
	mockServices = []Service{
		Service{
			ID:           1,
			Name:         "test",
			GroupID:      1,
			DatacenterID: 1,
		},
		Service{
			ID:           2,
			Name:         "test2",
			GroupID:      2,
			DatacenterID: 3,
		},
	}
)

func getServiceSubcriber() {
	n.Subscribe("service.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qs := Service{}
			json.Unmarshal(msg.Data, &qs)

			for _, service := range mockServices {
				if qs.GroupID != 0 && service.GroupID == qs.GroupID && service.Name == qs.Name {
					data, _ := json.Marshal(service)
					n.Publish(msg.Reply, data)
					return
				} else if qs.GroupID == 0 && service.Name == qs.Name {
					data, _ := json.Marshal(service)
					n.Publish(msg.Reply, data)
					return
				}
			}
		}
		n.Publish(msg.Reply, []byte(`{"error":"not found"}`))
	})
}

func findServiceSubcriber() {
	n.Subscribe("service.find", func(msg *nats.Msg) {
		data, _ := json.Marshal(mockServices)
		n.Publish(msg.Reply, data)
	})
}

func createServiceSubcriber() {
	n.Subscribe("service.set", func(msg *nats.Msg) {
		var s Service

		json.Unmarshal(msg.Data, &s)
		s.ID = 3
		data, _ := json.Marshal(s)

		n.Publish(msg.Reply, data)
	})
}

func deleteServiceSubcriber() {
	n.Subscribe("service.del", func(msg *nats.Msg) {
		var s Service

		json.Unmarshal(msg.Data, &s)

		n.Publish(msg.Reply, []byte{})
	})
}

func TestServices(t *testing.T) {
	Convey("Given service handler", t, func() {
		// setup nats connection
		os.Setenv("JWT_SECRET", "test")
		setup()

		Convey("When getting a list of services", func() {
			findServiceSubcriber()

			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
			c.SetPath("/services/")

			Convey("It should return the correct set of data", func() {
				var s []Service

				err := getServicesHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &s)

				So(err, ShouldBeNil)
				So(len(s), ShouldEqual, 2)
				So(s[0].ID, ShouldEqual, 1)
				So(s[0].Name, ShouldEqual, "test")
				So(s[0].GroupID, ShouldEqual, 1)
			})

		})

		Convey("When getting a single service", func() {
			getServiceSubcriber()

			Convey("Where the authenticated user is an admin", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "admin"
				ft.Claims["admin"] = true
				ft.Claims["group_id"] = 2.0

				c.SetPath("/services/:service")
				c.SetParamNames("service")
				c.SetParamValues("test")
				c.Set("user", ft)

				Convey("It should return the correct set of data", func() {
					var s Service

					err := getServiceHandler(c)
					So(err, ShouldBeNil)

					resp := rec.Body.Bytes()
					err = json.Unmarshal(resp, &s)

					So(err, ShouldBeNil)
					So(s.ID, ShouldEqual, 1)
					So(s.Name, ShouldEqual, "test")
				})
			})

			Convey("Where the service group matches the authenticated users group", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "admin"
				ft.Claims["admin"] = false
				ft.Claims["group_id"] = 1.0

				c.SetPath("/services/:service")
				c.SetParamNames("service")
				c.SetParamValues("test")
				c.Set("user", ft)

				Convey("It should return the correct set of data", func() {
					var s Service

					err := getServiceHandler(c)
					So(err, ShouldBeNil)

					resp := rec.Body.Bytes()
					err = json.Unmarshal(resp, &s)

					So(err, ShouldBeNil)
					So(s.ID, ShouldEqual, 1)
					So(s.Name, ShouldEqual, "test")
				})
			})

			Convey("Where the service group does not match the authenticated users group", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "test2"
				ft.Claims["admin"] = false
				ft.Claims["group_id"] = 2.0

				c.SetPath("/services/:service")
				c.SetParamNames("service")
				c.SetParamValues("test")
				c.Set("user", ft)

				Convey("It should return an 404 doesn't exist", func() {
					err := getServiceHandler(c)
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
				})
			})
		})

	})
}
