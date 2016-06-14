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
	mockGroups = []Group{
		Group{
			ID:   "1",
			Name: "test",
		},
		Group{
			ID:   "2",
			Name: "test",
		},
	}
)

func getGroupSubcriber() {
	n.Subscribe("group.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qg := Group{}
			json.Unmarshal(msg.Data, &qg)
			for _, group := range mockGroups {
				if group.ID == qg.ID || group.Name == qg.Name {
					data, _ := json.Marshal(group)
					n.Publish(msg.Reply, data)
					return
				}
			}
			n.Publish(msg.Reply, []byte(`{"error":"not found"}`))
		}

		data, _ := json.Marshal(mockGroups)
		n.Publish(msg.Reply, data)
	})
}

func createGroupSubcriber() {
	n.Subscribe("group.create", func(msg *nats.Msg) {
		var u Group

		json.Unmarshal(msg.Data, &u)
		u.ID = "3"
		data, _ := json.Marshal(u)

		n.Publish(msg.Reply, data)
	})
}

func deleteGroupSubcriber() {
	n.Subscribe("group.delete.1", func(msg *nats.Msg) {
		n.Publish(msg.Reply, []byte{})
	})
}

func TestGroups(t *testing.T) {
	Convey("Given group handler", t, func() {
		// setup nats connection
		os.Setenv("JWT_SECRET", "test")
		setup()

		Convey("When getting a list of groups", func() {
			getGroupSubcriber()

			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
			c.SetPath("/groups/")

			Convey("It should return the correct set of data", func() {
				var g []Group

				err := getGroupsHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &g)

				So(err, ShouldBeNil)
				So(len(g), ShouldEqual, 2)
				So(g[0].ID, ShouldEqual, "1")
				So(g[0].Name, ShouldEqual, "test")
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
			c.SetParamValues("1")

			Convey("It should return the correct set of data", func() {
				var g Group

				err := getGroupHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &g)

				So(err, ShouldBeNil)
				So(g.ID, ShouldEqual, "1")
				So(g.Name, ShouldEqual, "test")
			})

		})

		Convey("When creating a group", func() {
			createGroupSubcriber()

			data, _ := json.Marshal(Group{Name: "new-test"})

			e := echo.New()
			req, _ := http.NewRequest("POST", "/groups/", bytes.NewReader(data))
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			c.SetPath("/groups/")

			Convey("It should create the group and return the correct set of data", func() {
				var g Group

				err := createGroupHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &g)

				So(err, ShouldBeNil)
				So(g.ID, ShouldEqual, "3")
				So(g.Name, ShouldEqual, "new-test")
			})

		})

		Convey("When deleting a group", func() {
			deleteGroupSubcriber()

			e := echo.New()
			req := http.Request{Method: "DELETE"}
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(&req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			c.SetPath("/groups/:group")
			c.SetParamNames("group")
			c.SetParamValues("1")

			Convey("It should delete the group and return ok", func() {
				err := deleteGroupHandler(c)
				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})

		})

	})
}
