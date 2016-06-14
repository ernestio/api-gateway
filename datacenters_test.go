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
	mockDatacenters = []Datacenter{
		Datacenter{
			ID:      "1",
			Name:    "test",
			GroupID: "1",
		},
		Datacenter{
			ID:      "2",
			Name:    "test2",
			GroupID: "2",
		},
	}
)

func getDatacenterSubcriber() {
	n.Subscribe("datacenter.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qd := Datacenter{}
			json.Unmarshal(msg.Data, &qd)

			for _, datacenter := range mockDatacenters {
				if qd.GroupID != "" && datacenter.GroupID == qd.GroupID && datacenter.Name == qd.Name {
					data, _ := json.Marshal(datacenter)
					n.Publish(msg.Reply, data)
					return
				} else if qd.GroupID == "" && datacenter.Name == qd.Name {
					data, _ := json.Marshal(datacenter)
					n.Publish(msg.Reply, data)
					return
				}
			}
		}
		n.Publish(msg.Reply, []byte(`{"error":"not found"}`))
	})
}

func findDatacenterSubcriber() {
	n.Subscribe("datacenter.find", func(msg *nats.Msg) {
		data, _ := json.Marshal(mockDatacenters)
		n.Publish(msg.Reply, data)
	})
}

func createDatacenterSubcriber() {
	n.Subscribe("datacenter.set", func(msg *nats.Msg) {
		var d Datacenter

		json.Unmarshal(msg.Data, &d)
		d.ID = "3"
		data, _ := json.Marshal(d)

		n.Publish(msg.Reply, data)
	})
}

func deleteDatacenterSubcriber() {
	n.Subscribe("datacenter.del", func(msg *nats.Msg) {
		var u Datacenter

		json.Unmarshal(msg.Data, &u)

		n.Publish(msg.Reply, []byte{})
	})
}

func TestDatacenters(t *testing.T) {
	Convey("Given datacenter handler", t, func() {
		// setup nats connection
		os.Setenv("JWT_SECRET", "test")
		setup()

		Convey("When getting a list of datacenters", func() {
			findDatacenterSubcriber()

			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
			c.SetPath("/datacenters/")

			Convey("It should return the correct set of data", func() {
				var d []Datacenter

				err := getDatacentersHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &d)

				So(err, ShouldBeNil)
				So(len(d), ShouldEqual, 2)
				So(d[0].ID, ShouldEqual, "1")
				So(d[0].Name, ShouldEqual, "test")
			})

		})

		Convey("When getting a single datacenter", func() {
			getDatacenterSubcriber()

			Convey("Where the authenticated user is an admin", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "admin"
				ft.Claims["admin"] = true
				ft.Claims["group_id"] = "2"

				c.SetPath("/datacenters/:datacenter")
				c.SetParamNames("datacenter")
				c.SetParamValues("test")
				c.Set("user", ft)

				Convey("It should return the correct set of data", func() {
					var d Datacenter

					err := getDatacenterHandler(c)
					So(err, ShouldBeNil)

					resp := rec.Body.Bytes()
					err = json.Unmarshal(resp, &d)

					So(err, ShouldBeNil)
					So(d.ID, ShouldEqual, "1")
					So(d.Name, ShouldEqual, "test")
				})
			})

			Convey("Where the datacenter group matches the authenticated users group", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "admin"
				ft.Claims["admin"] = false
				ft.Claims["group_id"] = "1"

				c.SetPath("/datacenters/:datacenter")
				c.SetParamNames("datacenter")
				c.SetParamValues("test")
				c.Set("user", ft)

				Convey("It should return the correct set of data", func() {
					var d Datacenter

					err := getDatacenterHandler(c)
					So(err, ShouldBeNil)

					resp := rec.Body.Bytes()
					err = json.Unmarshal(resp, &d)

					So(err, ShouldBeNil)
					So(d.ID, ShouldEqual, "1")
					So(d.Name, ShouldEqual, "test")
				})
			})

			Convey("Where the datacenter group does not match the authenticated users group", func() {
				e := echo.New()
				req := new(http.Request)
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "test2"
				ft.Claims["admin"] = false
				ft.Claims["group_id"] = "2"

				c.SetPath("/datacenters/:datacenter")
				c.SetParamNames("datacenter")
				c.SetParamValues("test")
				c.Set("user", ft)

				Convey("It should return an 404 doesn't exist", func() {
					err := getDatacenterHandler(c)
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
				})
			})
		})

		Convey("When creating a datacenter", func() {
			createDatacenterSubcriber()

			mockDC := Datacenter{
				GroupID:   "1",
				Name:      "new-test",
				Type:      "vcloud",
				Username:  "test",
				Password:  "test",
				VCloudURL: "test",
			}

			data, _ := json.Marshal(mockDC)

			Convey("Where the authenticated user is an admin", func() {
				e := echo.New()
				req, _ := http.NewRequest("POST", "/datacenters/", bytes.NewReader(data))
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "admin"
				ft.Claims["admin"] = true
				ft.Claims["group_id"] = "2"

				c.SetPath("/datacenters/")
				c.Set("user", ft)

				Convey("It should create the datacenter and return the correct set of data", func() {
					var d Datacenter

					err := createDatacenterHandler(c)
					So(err, ShouldBeNil)

					resp := rec.Body.Bytes()
					err = json.Unmarshal(resp, &d)

					So(err, ShouldBeNil)
					So(d.ID, ShouldEqual, "3")
					So(d.Name, ShouldEqual, "new-test")
				})
			})

			Convey("Where the datacenter group matches the authenticated users group", func() {
				e := echo.New()
				req, _ := http.NewRequest("POST", "/datacenters/", bytes.NewReader(data))
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "test"
				ft.Claims["admin"] = true
				ft.Claims["group_id"] = "1"

				c.SetPath("/datacenters/")
				c.Set("user", ft)

				Convey("It should create the datacenter and return the correct set of data", func() {
					var d Datacenter

					err := createDatacenterHandler(c)
					So(err, ShouldBeNil)

					resp := rec.Body.Bytes()
					err = json.Unmarshal(resp, &d)

					So(err, ShouldBeNil)
					So(d.ID, ShouldEqual, "3")
					So(d.Name, ShouldEqual, "new-test")
				})
			})

			Convey("Where the datacenter group does not match the authenticated users group", func() {
				e := echo.New()
				req, _ := http.NewRequest("POST", "/datacenters/", bytes.NewReader(data))
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "admin"
				ft.Claims["admin"] = false
				ft.Claims["group_id"] = "2"

				c.SetPath("/datacenters/")
				c.Set("user", ft)

				Convey("It should return an 403 unauthorized error", func() {
					err := createDatacenterHandler(c)
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
				})
			})
		})

		Convey("When deleting a datacenter", func() {
			deleteDatacenterSubcriber()

			e := echo.New()
			req := http.Request{Method: "DELETE"}
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(&req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			ft := jwt.New(jwt.SigningMethodHS256)
			ft.Claims["username"] = "test"
			ft.Claims["admin"] = false
			ft.Claims["group_id"] = "1"

			c.SetPath("/datacenters/:datacenter")
			c.SetParamNames("datacenter")
			c.SetParamValues("test")
			c.Set("user", ft)

			Convey("It should delete the datacenter and return ok", func() {
				err := deleteDatacenterHandler(c)
				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})

		})

	})
}
