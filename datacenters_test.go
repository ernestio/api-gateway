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
	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mockDatacenters = []Datacenter{
		Datacenter{
			ID:   "1",
			Name: "test",
		},
		Datacenter{
			ID:   "2",
			Name: "test",
		},
	}
)

func getDatacenterSubcriber() {
	n.Subscribe("datacenter.get", func(msg *nats.Msg) {
		data, _ := json.Marshal(mockDatacenters[0])
		n.Publish(msg.Reply, data)
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
		var u Datacenter

		json.Unmarshal(msg.Data, &u)
		u.ID = "3"
		data, _ := json.Marshal(u)

		n.Publish(msg.Reply, data)
	})
}

func deleteDatacenterSubcriber() {
	n.Subscribe("datacenter.del", func(msg *nats.Msg) {
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

			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			c.SetPath("/datacenters/:datacenter")
			c.SetParamNames("datacenter")
			c.SetParamValues("1")

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

		Convey("When creating a datacenter", func() {
			createDatacenterSubcriber()

			data, _ := json.Marshal(Datacenter{Name: "new-test"})

			e := echo.New()
			req, _ := http.NewRequest("POST", "/datacenters/", bytes.NewReader(data))
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			c.SetPath("/datacenters/")

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

		Convey("When deleting a datacenter", func() {
			deleteDatacenterSubcriber()

			e := echo.New()
			req := http.Request{Method: "DELETE"}
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(&req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			c.SetPath("/datacenters/:datacenter")
			c.SetParamNames("datacenter")
			c.SetParamValues("1")

			Convey("It should delete the datacenter and return ok", func() {
				err := deleteDatacenterHandler(c)
				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})

		})

	})
}
