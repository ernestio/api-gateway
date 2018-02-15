/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/ernestio/api-gateway/controllers/envs"
	"github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListingEnvs(t *testing.T) {
	testsSetup()
	au := mockUsers[0]

	Convey("Scenario: getting a list of services", t, func() {
		Convey("Given services exist on the store", func() {
			Convey("When I call GET /services/", func() {
				foundSubscriber("environment.find", `[{"id":1,"name":"fake/test"},{"id":2,"name":"fake/test"}]`, 1)
				au.Admin = helpers.Bool(true)
				s, b := envs.List(au, nil)
				Convey("It should return the correct set of data", func() {
					var sr []models.Env
					So(s, ShouldEqual, 200)
					err := json.Unmarshal(b, &sr)
					So(err, ShouldBeNil)
					So(len(sr), ShouldEqual, 2)
					So(sr[0].ID, ShouldEqual, 1)
					So(sr[0].Name, ShouldEqual, "test")
				})
			})
		})
	})
}

func TestGetEnv(t *testing.T) {
	testsSetup()
	au := mockUsers[0]

	Convey("Scenario: getting a single service", t, func() {
		Convey("Given the service do not exist on the store", func() {
			foundSubscriber("environment.get", `{"_error":"not found"}`, 1)
			foundSubscriber("authorization.find", `[{"role":"reader"}]`, 1)
			Convey("And I call /service/:service on the api", func() {
				s, _ := envs.Get(au, "1")
				So(s, ShouldEqual, 404)
			})
		})
	})
}

func TestSearchEnv(t *testing.T) {
	testsSetup()
	au := mockUsers[0]

	Convey("Scenario: searching for services", t, func() {
		Convey("Given the service doesn't exist on the store", func() {
			foundSubscriber("environment.find", `[]`, 1)
			foundSubscriber("authorization.find", `["role":"owner"]`, 1)
			Convey("And I call /service/search/ on the api", func() {
				var s []views.BuildRender
				params := make(map[string]interface{})
				params["service"] = "1"
				st, resp := envs.Search(au, params)

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should return an empty array", func() {
						err := json.Unmarshal(resp, &s)
						So(err, ShouldBeNil)
						So(st, ShouldEqual, 200)
						So(len(s), ShouldEqual, 0)
					})
				})
			})
		})
	})
}
