/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/ernestio/api-gateway/config"
	"github.com/ernestio/api-gateway/controllers/builds"
	"github.com/ernestio/api-gateway/controllers/envs"
	"github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResetEnvironment(t *testing.T) {
	testsSetup()
	config.Setup()
	au := mockUsers[0]

	Convey("Scenario: reseting a service", t, func() {
		action := models.Action{Type: "reset"}
		findUserSubscriber()
		foundSubscriber("environment.set", `"success"`, 1)

		Convey("Given my existing service is in progress", func() {
			foundSubscriber("environment.get", `{"id":1,"name":"fake/test","status":"in_progress"}`, 1)
			foundSubscriber("build.find", `[{"id":"1","name":"fake/test","status":"in_progress"}]`, 1)
			foundSubscriber("authorization.find", `[{"role":"owner"}]`, 1)
			serviceResetSubscriber()
			Convey("When I do a call to /services/reset", func() {
				s, b := envs.Reset(au, "fake/test", &action)
				Convey("Then it should return a success message", func() {
					So(s, ShouldEqual, 200)
					So(string(b), ShouldEqual, `success`)
				})
			})
		})

		Convey("Given my existing service is errored", func() {
			foundSubscriber("environment.get", `{"id":1,"name":"fake/test","status":"in_progress"}`, 1)
			foundSubscriber("build.find", `[{"id":"1","name":"fake/test","status":"errored"}]`, 1)
			foundSubscriber("authorization.find", `[{"role":"owner","resource_id":"1"}]`, 1)

			Convey("When I do a call to /services/reset", func() {
				s, b := envs.Reset(au, "fake/test", &action)
				Convey("Then it should return an error message", func() {
					So(s, ShouldEqual, 200)
					So(string(b), ShouldEqual, "Reset only applies to an 'in progress' environment, however environment 'fake/test' is on status 'errored")
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
	config.Setup()
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
	config.Setup()
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

func TestDeletingEnvs(t *testing.T) {
	testsSetup()
	config.Setup()
	au := mockUsers[0]

	Convey("Scenario: deleting a service", t, func() {
		Convey("Given a service exists with in progress status", func() {
			foundSubscriber("environment.get", `{"id":1,"status":"in_progress"}`, 3)
			foundSubscriber("build.find", `[{"id":"test","status":"in_progress"}]`, 1)
			foundSubscriber("build.get.mapping", `{}`, 1)
			foundSubscriber("build.set", `{"_error": "environment build is in progress"}`, 1)
			foundSubscriber("mapping.get.delete", `{"id":"test-uuid-1"}`, 1)
			foundSubscriber("datacenter.get", `{"id":1, "credentials": {"username":" test"}}`, 1)
			res := `[{"resource_id":"1","role":"reader"}]`
			foundSubscriber("authorization.find", res, 1)
			Convey("When I call DELETE /services/:service", func() {
				st, resp := builds.Delete(au, "foo-bar")
				Convey("Then I should get a 400 response", func() {
					So(st, ShouldEqual, 400)
					So(string(resp), ShouldEqual, `"Environment is already applying some changes, please wait until they are done"`)
				})
			})
		})
		Convey("Given a service exists on the store", func() {
			foundSubscriber("environment.get", `{"id":1,"status":"done"}`, 3)
			foundSubscriber("build.find", `[{"id":"test","status":"done"}]`, 1)
			foundSubscriber("build.get.mapping", `{}`, 1)
			foundSubscriber("build.set", `{}`, 1)
			foundSubscriber("mapping.get.delete", `{"id":"foo-bar"}`, 1)
			foundSubscriber("datacenter.get", `{"id":1, "credentials": {"username":" test"}}`, 1)
			res := `[{"resource_id":"1","role":"reader"}]`
			foundSubscriber("authorization.find", res, 1)
			Convey("When I call DELETE /services/:service", func() {
				st, resp := builds.Delete(au, "foo-bar")

				Convey("Then I should get a response with id and stream id", func() {
					So(st, ShouldEqual, 200)
					So(string(resp), ShouldEqual, `{"id":"foo-bar"}`)
				})
			})
		})
	})
}
