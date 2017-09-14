/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/ernestio/api-gateway/config"
	"github.com/ernestio/api-gateway/controllers/envs"
	"github.com/ernestio/api-gateway/views"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResetEnvironment(t *testing.T) {
	testsSetup()
	config.Setup()
	au := mockUsers[0]

	Convey("Scenario: reseting a service", t, func() {
		findUserSubscriber()
		foundSubscriber("service.set", `"success"`, 1)

		Convey("Given my existing service is in progress", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"fake/test","status":"in_progress"},{"id":"2","name":"fake/test","status":"done"}]`, 1)
			foundSubscriber("authorization.find", `[{"role":"owner"}]`, 1)
			serviceResetSubscriber()
			Convey("When I do a call to /services/reset", func() {
				s, b := envs.Reset(au, "foo")
				Convey("Then it should return a success message", func() {
					So(s, ShouldEqual, 200)
					So(string(b), ShouldEqual, `success`)
				})
			})
		})

		Convey("Given my existing service is errored", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"fake/test","status":"errored"},{"id":"2","name":"fake/test","status":"done"}]`, 1)
			foundSubscriber("authorization.find", `[{"role":"owner","resource_id":"1"}]`, 1)

			Convey("When I do a call to /services/reset", func() {
				s, b := envs.Reset(au, "foo")
				Convey("Then it should return an error message", func() {
					So(s, ShouldEqual, 200)
					So(string(b), ShouldEqual, "Reset only applies to an 'in progress' environment, however environment 'foo' is on status 'errored")
				})
			})
		})
	})
}

func TestListingEnvs(t *testing.T) {
	testsSetup()
	config.Setup()
	au := mockUsers[0]

	Convey("Scenario: getting a list of services", t, func() {
		Convey("Given services exist on the store", func() {
			foundSubscriber("user.find", `[{"id":"1"}]`, 1)
			foundSubscriber("service.find", `[{"id":"1","name":"fake/test","datacenter_id":1},{"id":"2","name":"fake/test","datacenter_id":2}]`, 1)
			Convey("When I call GET /services/", func() {
				au.Admin = true
				s, b := envs.List(au)
				Convey("It should return the correct set of data", func() {
					var sr []views.ServiceRender
					So(s, ShouldEqual, 200)
					err := json.Unmarshal(b, &sr)
					So(err, ShouldBeNil)
					So(len(sr), ShouldEqual, 1)
					So(sr[0].ID, ShouldEqual, "1")
					So(sr[0].Name, ShouldEqual, "test")
					So(sr[0].DatacenterID, ShouldEqual, 1)
				})
			})
		})
	})
}

func TestGetEnv(t *testing.T) {
	testsSetup()
	config.Setup()
	au := mockUsers[0]

	Convey("Scenario: getting a single service", t, func() {
		Convey("Given the service do not exist on the store", func() {
			foundSubscriber("service.find", `[]`, 1)
			Convey("And I call /service/:service on the api", func() {
				s, _ := envs.Get(au, "1")
				So(s, ShouldEqual, 404)
			})
		})
	})
}

func TestSearchEnv(t *testing.T) {
	testsSetup()
	config.Setup()
	au := mockUsers[0]

	Convey("Scenario: searching for services", t, func() {
		Convey("Given the service doesn't exist on the store", func() {
			foundSubscriber("service.find", `[]`, 1)
			foundSubscriber("authorization.find", `["role":"owner"]`, 1)
			Convey("And I call /service/search/ on the api", func() {
				var s []views.ServiceRender
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

func TestDeletingEnvs(t *testing.T) {
	testsSetup()
	config.Setup()
	au := mockUsers[0]

	Convey("Scenario: deleting a service", t, func() {
		Convey("Given a service exists with in progress status", func() {
			foundSubscriber("service.find", `[{"id":"foo-bar","status":"in_progress"}]`, 1)
			res := `[{"resource_id":"1","role":"owner"}]`
			foundSubscriber("authorization.find", res, 1)
			Convey("When I call DELETE /services/:service", func() {
				st, resp := envs.Delete(au, "foo-bar")
				SkipConvey("Then I should get a 400 response", func() {
					So(st, ShouldEqual, 400)
					So(string(resp), ShouldEqual, `"Environment is already applying some changes, please wait until they are done"`)
				})
			})
		})
		Convey("Given a service exists on the store", func() {
			foundSubscriber("service.find", `[{"id":"foo-bar","status":"done"}]`, 1)
			foundSubscriber("definition.map.deletion", `{}`, 1)
			foundSubscriber("service.delete", `"success"`, 1)
			res := `[{"resource_id":"1","role":"owner"}]`
			foundSubscriber("authorization.find", res, 1)
			Convey("When I call DELETE /services/:service", func() {
				st, resp := envs.Delete(au, "foo-bar")

				Convey("Then I should get a response with id and stream id", func() {
					So(st, ShouldEqual, 200)
					So(string(resp), ShouldEqual, `{"id":"foo-bar"}`)
				})
			})
		})
	})
}
