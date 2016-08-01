/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServices(t *testing.T) {
	testsSetup()
	setup()

	Convey("Scenario: reeting a service", t, func() {
		foundSubscriber("service.set", `"success"`, 1)

		Convey("Given my existing service is in progress", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"test","status":"in_progress"},{"id":"2","name":"test","status":"done"}]`, 1)

			Convey("When I do a call to /services/reset", func() {
				params := make(map[string]string)
				params["service"] = "foo"
				resp, err := doRequest("POST", "/services/foo/reset/", params, nil, resetServiceHandler, nil)
				Convey("Then it should return a success message", func() {
					So(err, ShouldBeNil)
					So(string(resp), ShouldEqual, `success`)
				})
			})
		})

		Convey("Given my existing service is errored", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"test","status":"errored"},{"id":"2","name":"test","status":"done"}]`, 1)

			Convey("When I do a call to /services/reset", func() {
				params := make(map[string]string)
				params["service"] = "foo"
				resp, err := doRequest("POST", "/services/foo/reset/", params, nil, resetServiceHandler, nil)
				Convey("Then it should return an error message", func() {
					So(err, ShouldBeNil)
					So(string(resp), ShouldEqual, "Reset only applies to 'in progress' serices, however service 'foo' is on status 'errored")
				})
			})
		})
	})

	Convey("Scenario: generating a uuid", t, func() {
		Convey("Given I do a call to /services/uuid", func() {
			resp, err := doRequest("POST", "/services/uuid/", nil, []byte(`{"id":"foo"}`), createUUIDHandler, nil)

			Convey("It should return the correct encoded uuid", func() {
				So(err, ShouldBeNil)
				So(string(resp), ShouldEqual, `{"uuid":"acbd18db4cc2f85cedef654fccc4a4d8"}`)
			})
		})
	})

	Convey("Scenario: getting a list of services", t, func() {
		Convey("Given services exist on the store", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"test","datacenter_id":1},{"id":"2","name":"test","datacenter_id":2}]`, 2)
			Convey("When I call GET /services/", func() {
				resp, err := doRequest("GET", "/services/", nil, nil, getServicesHandler, nil)

				Convey("It should return the correct set of data", func() {
					var s []OutputService
					So(err, ShouldBeNil)
					err = json.Unmarshal(resp, &s)
					So(err, ShouldBeNil)
					So(len(s), ShouldEqual, 2)
					So(s[0].ID, ShouldEqual, "1")
					So(s[0].Name, ShouldEqual, "test")
					So(s[0].DatacenterID, ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Scenario: getting a single service", t, func() {
		Convey("Given the service do not exist on the store", func() {
			foundSubscriber("service.find", `[]`, 2)
			Convey("And I call /service/:service on the api", func() {
				params := make(map[string]string)
				params["service"] = "1"
				resp, err := doRequest("GET", "/services/:service", params, nil, getServiceHandler, nil)
				So(string(resp), ShouldEqual, "null")
				So(err, ShouldBeNil)
			})
		})
		Convey("Given the service exists on the store", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"test","datacenter_id":1},{"id":"2","name":"test","datacenter_id":2}]`, 2)
			Convey("And I call /service/:service on the api", func() {
				var d OutputService
				params := make(map[string]string)
				params["service"] = "1"
				resp, err := doRequest("GET", "/services/:service", params, nil, getServiceHandler, nil)

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should get the existing service", func() {

						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &d)

						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, "1")
						So(d.Name, ShouldEqual, "test")
					})
				})

				Convey("When the service group matches the authenticated users group", func() {
					ft := generateTestToken(1, "test", false)

					params := make(map[string]string)
					params["service"] = "1"
					resp, err := doRequest("GET", "/services/:service", params, nil, getServiceHandler, ft)

					Convey("Then I should get the existing service", func() {
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, "1")
						So(d.Name, ShouldEqual, "test")
					})
				})
			})
		})
	})

	Convey("Scenario: getting a service's builds", t, func() {
		Convey("Given the service exists on the store", func() {
			Convey("And I call /service/:service/builds/ on the api", func() {
				findServiceSubscriber()
				var s []OutputService
				params := make(map[string]string)
				params["service"] = "test"
				resp, err := doRequest("GET", "/services/:service/builds/", params, nil, getServiceBuildsHandler, nil)

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should get the service's builds", func() {
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &s)

						So(err, ShouldBeNil)
						So(len(s), ShouldEqual, 2)
						So(s[0].ID, ShouldEqual, "1")
						So(s[0].Name, ShouldEqual, "test")
					})
				})

				Convey("When the service group matches the authenticated users group", func() {
					findServiceSubscriber()
					ft := generateTestToken(1, "test", false)

					params := make(map[string]string)
					params["service"] = "test"
					resp, err := doRequest("GET", "/services/:service/builds/", params, nil, getServiceBuildsHandler, ft)

					Convey("Then I should get the service's builds", func() {
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &s)

						So(len(s), ShouldEqual, 2)
						So(s[0].ID, ShouldEqual, "1")
						So(s[0].Name, ShouldEqual, "test")
					})
				})
			})
		})
	})

	Convey("Scenario: getting a service's build", t, func() {
		Convey("Given the service exists on the store", func() {
			Convey("And I call /service/:service/builds/:build on the api", func() {
				findServiceSubscriber()
				var s OutputService

				params := make(map[string]string)
				params["service"] = "test"
				params["id"] = "1"
				resp, err := doRequest("GET", "/services/:service/builds/:build", params, nil, getServiceBuildHandler, nil)

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should get the existing service", func() {
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &s)

						So(err, ShouldBeNil)
						So(s.ID, ShouldEqual, "1")
						So(s.Name, ShouldEqual, "test")
					})
				})

				Convey("When the service group matches the authenticated users group", func() {
					findServiceSubscriber()
					ft := generateTestToken(1, "test", false)

					params := make(map[string]string)
					params["service"] = "test"
					params["id"] = "1"
					resp, err := doRequest("GET", "/services/:service/builds/:build", params, nil, getServiceBuildHandler, ft)

					Convey("Then I should get the existing service", func() {
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &s)
						So(err, ShouldBeNil)
						So(s.ID, ShouldEqual, "1")
						So(s.Name, ShouldEqual, "test")
					})
				})
			})
		})
	})

	Convey("Scenario: searching for services", t, func() {
		Convey("Given the service exists on the store", func() {
			findServiceSubscriber()
			Convey("And I call /service/search/ on the api", func() {
				var s []OutputService
				params := make(map[string]string)
				params["service"] = "1"
				resp, err := doRequest("GET", "/services/search/?name=test", params, nil, searchServicesHandler, nil)

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should get the matching service", func() {
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &s)

						So(err, ShouldBeNil)
						So(len(s), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("Given the service doesn't exist on the store", func() {
			findServiceSubscriber()
			Convey("And I call /service/search/ on the api", func() {
				var s []OutputService
				params := make(map[string]string)
				params["service"] = "1"
				resp, err := doRequest("GET", "/services/search/?name=doesntexist", params, nil, searchServicesHandler, nil)

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should return an empty array", func() {
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &s)

						So(err, ShouldBeNil)
						So(len(s), ShouldEqual, 0)
					})
				})
			})
		})
	})

	Convey("Scenario: creating a service", t, func() {
		params := make(map[string]string)
		Convey("Given I do a post call to /services ", func() {
			createServiceSubscriber()

			Convey("And the content type is non json and non yaml", func() {
				data := []byte("bla")
				resp, err := doRequest("POST", "/services/", params, data, createServiceHandler, nil)
				Convey("Then I should get a 400 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Invalid input format"`)
				})
			})

			Convey("And the content type is a non valid yaml", func() {
				data := []byte("asd")
				headers := map[string]string{}
				headers["Content-Type"] = "application/yaml"
				resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
				Convey("Then I should get a 400 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Invalid input"`)
				})
			})

			Convey("And the content type is a non valid json", func() {
				data := []byte(`{"name"}`)
				headers := map[string]string{}
				headers["Content-Type"] = "application/json"
				resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
				Convey("Then I should get a 400 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Invalid input"`)
				})
			})

			Convey("And the specified datacenter does not exist", func() {
				foundSubscriber("datacenter.find", "[]", 1)
				data := []byte(`{"name":"test"}`)
				headers := map[string]string{}
				headers["Content-Type"] = "application/json"
				resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
				SkipConvey("Then I should get a 404 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Specified datacenter does not exist"`)
				})
			})

			Convey("And the specified group does not exist", func() {
				notFoundSubscriber("group.get", 1)
				foundSubscriber("datacenter.find", `[{"id":"1"}]`, 1)
				data := []byte(`{"name":"test"}`)
				headers := map[string]string{}
				headers["Content-Type"] = "application/json"
				resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
				Convey("Then I should get a 404 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Specified group does not exist"`)
				})
			})

			Convey("And I provide a valid input", func() {
				foundSubscriber("group.get", `{"id":"1"}`, 1)
				foundSubscriber("datacenter.find", `[{"id":"1"}]`, 1)
				data := []byte(`{"name":"test"}`)
				headers := map[string]string{}
				headers["Content-Type"] = "application/json"

				Convey("And the service does not exist", func() {
					foundSubscriber("service.find", "[]", 1)
					foundSubscriber("service.create", `{"id":"1"}`, 1)
					foundSubscriber("definition.map.creation", `{"id":"1"}`, 1)
					resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
					Convey("Then I should get a response with a valid id", func() {
						So(err, ShouldEqual, nil)
						So(strings.Contains(string(resp), `{"id":"`), ShouldEqual, true)
						So(strings.Contains(string(resp), `-d29d2764b65cae3f4114164bb6cf80cb`), ShouldEqual, true)
					})
				})

				Convey("And the service already exists", func() {
					foundSubscriber("service.create", `{"id":"1"}`, 1)
					Convey("And the existing service is done", func() {
						foundSubscriber("definition.map.creation", `{"id":"1"}`, 1)
						foundSubscriber("service.find", `[{"id":"foo-bar","status":"done"}]`, 1)
						resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
						Convey("Then I should get a response with the existing id", func() {
							So(err, ShouldEqual, nil)
							So(strings.Contains(string(resp), `{"id":"`), ShouldEqual, true)
							So(strings.Contains(string(resp), `-d29d2764b65cae3f4114164bb6cf80cb`), ShouldEqual, true)
						})
					})

					Convey("And the existing service is in progress", func() {
						foundSubscriber("service.find", `[{"id":"foo-bar","status":"in_progress"}]`, 1)
						resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
						Convey("Then I should get an error as an in_progress service can't be modified", func() {
							So(err, ShouldEqual, nil)
							So(string(resp), ShouldEqual, `"Your service process is 'in progress' if your're sure you want to fix it please reset it first"`)
						})
					})

					Convey("And the existing service is errored", func() {
						foundSubscriber("service.find", `[{"id":"foo-bar","status":"errored"}]`, 1)
						foundSubscriber("definition.map.creation", `{"id":"1"}`, 1)
						foundSubscriber("service.patch", `{"id":"1"}`, 1)
						resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
						Convey("Then I should get a response with the existing id", func() {
							So(err, ShouldEqual, nil)
							So(strings.Contains(string(resp), `{"id":"`), ShouldEqual, true)
							So(strings.Contains(string(resp), `-d29d2764b65cae3f4114164bb6cf80cb`), ShouldEqual, true)
						})
					})
				})
			})
		})
	})

	Convey("Scenario: deleting a service", t, func() {
		ft := generateTestToken(1, "test", false)
		params := make(map[string]string)
		params["service"] = "foo-bar"

		Convey("Given I don't have services on my store", func() {
			foundSubscriber("service.find", `[]`, 1)
			Convey("When I call DELETE /services/:service", func() {
				_, err := doRequest("DELETE", "/services/:service", params, nil, deleteServiceHandler, ft)
				Convey("Then I should get a 400 response", func() {
					So(err.Error(), ShouldEqual, `"Service not found"`)
				})
			})
		})

		Convey("Given a service exists with in progress status", func() {
			foundSubscriber("service.find", `[{"id":"foo-bar","status":"in_progress"}]`, 1)
			Convey("When I call DELETE /services/:service", func() {
				req, err := doRequest("DELETE", "/services/:service", params, nil, deleteServiceHandler, ft)
				Convey("Then I should get a 400 response", func() {
					So(err, ShouldBeNil)
					So(string(req), ShouldEqual, `"Service is already applying some changes, please wait until they are done"`)
				})
			})
		})

		Convey("Given a service exists on the store", func() {
			foundSubscriber("service.find", `[{"id":"foo-bar","status":"done"}]`, 1)
			foundSubscriber("definition.map.deletion", `""`, 1)
			foundSubscriber("service.delete", `""`, 1)
			Convey("When I call DELETE /services/:service", func() {
				res, err := doRequest("DELETE", "/services/:service", params, nil, deleteServiceHandler, ft)

				Convey("Then I should get a response with id and stream id", func() {
					So(err, ShouldBeNil)
					So(string(res), ShouldEqual, `{"id":"foo-bar","stream_id":"bar"}`)
				})
			})

		})

	})
}
