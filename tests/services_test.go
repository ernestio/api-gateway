/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ernestio/api-gateway/config"
	"github.com/ernestio/api-gateway/controllers/services"
	"github.com/ernestio/api-gateway/views"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServices(t *testing.T) {
	testsSetup()
	config.Setup()
	au := mockUsers[0]

	Convey("Scenario: reseting a service", t, func() {
		findUserSubscriber()
		foundSubscriber("service.set", `"success"`, 1)

		Convey("Given my existing service is in progress", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"fake/test","status":"in_progress"},{"id":"2","name":"fake/test","status":"done"}]`, 1)
			serviceResetSubscriber()
			Convey("When I do a call to /services/reset", func() {
				s, b := services.Reset(au, "foo")
				Convey("Then it should return a success message", func() {
					So(s, ShouldEqual, 200)
					So(string(b), ShouldEqual, `success`)
				})
			})
		})

		Convey("Given my existing service is errored", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"fake/test","status":"errored"},{"id":"2","name":"fake/test","status":"done"}]`, 1)

			Convey("When I do a call to /services/reset", func() {
				s, b := services.Reset(au, "foo")
				Convey("Then it should return an error message", func() {
					So(s, ShouldEqual, 200)
					So(string(b), ShouldEqual, "Reset only applies to 'in progress' serices, however service 'foo' is on status 'errored")
				})
			})
		})
	})

	Convey("Scenario: getting a list of services", t, func() {
		Convey("Given services exist on the store", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"fake/test","datacenter_id":1},{"id":"2","name":"fake/test","datacenter_id":2}]`, 1)
			Convey("When I call GET /services/", func() {
				au.Admin = true
				s, b := services.List(au)
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

	Convey("Scenario: getting a single service", t, func() {
		Convey("Given the service do not exist on the store", func() {
			foundSubscriber("service.find", `[]`, 1)
			Convey("And I call /service/:service on the api", func() {
				s, _ := services.Get(au, "1")
				So(s, ShouldEqual, 404)
			})
		})

		Convey("Given the service exists on the store", func() {
			foundSubscriber("service.find", `[{"id":"1","name":"fake/test","datacenter_id":1},{"id":"2","name":"fake/test","datacenter_id":2}]`, 1)
			foundSubscriber("service.get.mapping", `{"name":"fake/test", "vpcs": {"items":[{"vpc_id":"22"}]}, "networks":{"items":[{"name":"a"}]}}`, 1)
			Convey("And I call /service/:service on the api", func() {
				var d views.ServiceRender
				Convey("When I'm authenticated as an admin user", func() {
					s, resp := services.Get(au, "1")
					Convey("Then I should get the existing service", func() {

						So(s, ShouldEqual, 200)
						err := json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, "1")
						So(d.Name, ShouldEqual, "test")
					})
				})
			})
		})
	})

	Convey("Scenario: searching for services", t, func() {
		Convey("Given the service exists on the store", func() {
			findServiceSubscriber()
			foundSubscriber("service.find", `[{"id":"1","name":"fake/test","datacenter_id":1},{"id":"2","name":"fake/test","datacenter_id":2}]`, 1)
			foundSubscriber("service.get.mapping", `{"name":"fake/test", "networks":{"items":[{"name":"a"}]}}`, 2)
			Convey("And I call /service/search/ on the api", func() {
				var s []views.ServiceRender
				params := make(map[string]interface{})
				params["name"] = "fake/test"

				Convey("When I'm authenticated as an admin user", func() {
					st, resp := services.Search(au, params)
					Convey("Then I should get the matching service", func() {
						err := json.Unmarshal(resp, &s)
						So(err, ShouldBeNil)
						So(st, ShouldEqual, 200)
						So(len(s), ShouldEqual, 2)
					})
				})
			})
		})
		Convey("Given the service doesn't exist on the store", func() {
			foundSubscriber("service.find", `[]`, 1)
			Convey("And I call /service/search/ on the api", func() {
				var s []views.ServiceRender
				params := make(map[string]interface{})
				params["service"] = "1"
				st, resp := services.Search(au, params)

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

	/*
		Convey("Scenario: creating a service", t, func() {
			params := make(map[string]string)
			Convey("Given I do a post call to /services ", func() {
					createServiceSubscriber()

					Convey("And the specified datacenter does not exist", func() {
						getUserSubscriber(1)
						foundSubscriber("datacenter.find", "[]", 1)
						data := []byte(`{"name":"test"}`)
						headers := map[string]string{}
						headers["Content-Type"] = "application/json"
						resp, err := doRequestHeaders("POST", "/services/", params, data, controllers.CreateServiceHandler, nil, headers)
						st, resp := services.Create(au, params)

						Convey("Then I should get a 404 response", func() {
							So(err, ShouldEqual, nil)
							So(string(resp), ShouldEqual, `"Specified datacenter does not exist"`)
						})
					})

						Convey("And the specified group does not exist", func() {
							notFoundSubscriber("group.get", 1)
							foundSubscriber("datacenter.find", `[{"id":1}]`, 1)
							data := []byte(`{"name":"test"}`)
							headers := map[string]string{}
							headers["Content-Type"] = "application/json"
							resp, _ := doRequestHeaders("POST", "/services/", params, data, controllers.CreateServiceHandler, nil, headers)
							Convey("Then I should get a 404 response", func() {
								// So(err.Error(), ShouldEqual, `code=404, message=`)
								// So(err, ShouldEqual, nil)
								So(string(resp), ShouldEqual, `"Specified group does not exist"`)
							})
						})

						Convey("And I provide a valid input", func() {
							foundSubscriber("group.get", `{"id":1}`, 1)
							foundSubscriber("datacenter.find", `[{"id":1}]`, 1)
							data := []byte(`{"name":"test"}`)
							headers := map[string]string{}
							headers["Content-Type"] = "application/json"

							SkipConvey("And the service does not exist", func() {
								foundSubscriber("service.find", "[]", 1)
								foundSubscriber("service.create", `{"id":"1"}`, 1)
								foundSubscriber("definition.map.creation", `{"id":"1"}`, 1)
								resp, err := doRequestHeaders("POST", "/services/", params, data, controllers.CreateServiceHandler, nil, headers)
								Convey("Then I should get a response with a valid id", func() {
									So(err, ShouldBeNil)
									So(strings.Contains(string(resp), `{"id":"`), ShouldEqual, true)
									So(strings.Contains(string(resp), `-d29d2764b65cae3f4114164bb6cf80cb`), ShouldEqual, true)
								})
							})

							SkipConvey("And the service already exists", func() {
								foundSubscriber("service.create", `{"id":"1"}`, 1)
								Convey("And the existing service is done", func() {
									foundSubscriber("definition.map.creation", `{"id":"1"}`, 1)
									foundSubscriber("service.find", `[{"id":"foo-bar","status":"done"}]`, 1)
									resp, err := doRequestHeaders("POST", "/services/", params, data, controllers.CreateServiceHandler, nil, headers)
									Convey("Then I should get a response with the existing id", func() {
										So(err, ShouldBeNil)
										So(strings.Contains(string(resp), `{"id":"`), ShouldEqual, true)
										So(strings.Contains(string(resp), `-d29d2764b65cae3f4114164bb6cf80cb`), ShouldEqual, true)
									})
								})

								Convey("And the existing service is in progress", func() {
									foundSubscriber("service.find", `[{"id":"foo-bar","status":"in_progress"}]`, 1)
									resp, err := doRequestHeaders("POST", "/services/", params, data, controllers.CreateServiceHandler, nil, headers)
									Convey("Then I should get an error as an in_progress service can't be modified", func() {
										So(err, ShouldEqual, nil)
										So(string(resp), ShouldEqual, `"Your service process is 'in progress' if your're sure you want to fix it please reset it first"`)
									})
								})

								Convey("And the existing service is errored", func() {
									foundSubscriber("service.find", `[{"id":"foo-bar","status":"errored"}]`, 1)
									foundSubscriber("definition.map.creation", `{"id":"1"}`, 1)
									foundSubscriber("service.patch", `{"id":"1"}`, 1)
									foundSubscriber("service.get.mapping", `{"id":"1"}`, 1)
									resp, err := doRequestHeaders("POST", "/services/", params, data, controllers.CreateServiceHandler, nil, headers)
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
	*/

	Convey("Scenario: deleting a service", t, func() {
		Convey("Given a service exists with in progress status", func() {
			foundSubscriber("service.find", `[{"id":"foo-bar","status":"in_progress"}]`, 1)
			res := `[{"resource_id":"1","role":"owner"}]`
			foundSubscriber("authorization.find", res, 1)
			Convey("When I call DELETE /services/:service", func() {
				st, resp := services.Delete(au, "foo-bar")
				Convey("Then I should get a 400 response", func() {
					So(st, ShouldEqual, 400)
					So(string(resp), ShouldEqual, `"Service is already applying some changes, please wait until they are done"`)
				})
			})
		})

		Convey("Given a service exists on the store", func() {
			foundSubscriber("service.find", `[{"id":"foo-bar","status":"done"}]`, 1)
			foundSubscriber("definition.map.deletion", `{}`, 1)
			foundSubscriber("service.delete", `""`, 1)
			res := `[{"resource_id":"1","role":"owner"}]`
			foundSubscriber("authorization.find", res, 1)
			Convey("When I call DELETE /services/:service", func() {
				st, resp := services.Delete(au, "foo-bar")

				Convey("Then I should get a response with id and stream id", func() {
					fmt.Println(string(resp))
					So(st, ShouldEqual, 200)
					So(string(resp), ShouldEqual, `{"id":"foo-bar"}`)
				})
			})

		})
	})
}
