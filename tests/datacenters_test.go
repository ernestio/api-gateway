/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/ernestio/api-gateway/config"
	"github.com/ernestio/api-gateway/controllers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDatacenters(t *testing.T) {
	testsSetup()
	config.Setup()

	Convey("Scenario: getting a list of datacenters", t, func() {
		Convey("Given datacenters exist on the store", func() {
			findDatacenterSubscriber()
			Convey("When I call /datacenters/", func() {
				resp, err := doRequest("GET", "/datacenters/", nil, nil, controllers.GetDatacentersHandler, nil)
				Convey("Then I should have a response with existing datacenters", func() {
					var d []models.Datacenter
					So(err, ShouldBeNil)

					err = json.Unmarshal(resp, &d)

					So(err, ShouldBeNil)
					So(len(d), ShouldEqual, 2)
					So(d[0].ID, ShouldEqual, 1)
					So(d[0].Name, ShouldEqual, "test")
				})
			})

			SkipConvey("Given no datacenters on the store", func() {
			})
		})
	})

	Convey("Scenario: getting a single datacenters", t, func() {
		Convey("Given the datacenter exists on the store", func() {
			getDatacenterSubscriber(2)
			Convey("And I call /datacenter/:datacenter on the api", func() {
				params := make(map[string]string)
				params["datacenter"] = "1"
				resp, err := doRequest("GET", "/datacenters/:datacenter", params, nil, controllers.GetDatacenterHandler, nil)

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should get the existing datacenter", func() {
						var d models.Datacenter

						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &d)

						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, 1)
						So(d.Name, ShouldEqual, "test")
					})
				})

				SkipConvey("When the datacenter group matches the authenticated users group", func() {
					ft := generateTestToken(1, "admin", true)

					params := make(map[string]string)
					params["datacenter"] = "1"
					resp, err := doRequest("GET", "/datacenters/:datacenter", params, nil, controllers.GetDatacenterHandler, ft)

					Convey("Then I should get the existing datacenter", func() {
						var d models.Datacenter
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, 1)
						So(d.Name, ShouldEqual, "test")
					})
				})

				SkipConvey("When the datacenter group does not match the authenticated users group", func() {
					ft := generateTestToken(2, "test2", false)
					params := make(map[string]string)
					params["datacenter"] = "1"
					_, err := doRequest("GET", "/datacenters/:datacenter", params, nil, controllers.GetDatacenterHandler, ft)

					Convey("Then I should get a 404 error as it doesn't exist", func() {
						So(err, ShouldNotBeNil)
						So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
					})
				})
			})
		})
	})

	Convey("Scenario: creating a datacenter", t, func() {
		Convey("Given the datacenter does not exist on the store ", func() {
			createDatacenterSubscriber()

			mockDC := models.Datacenter{
				GroupID:   1,
				Name:      "new-test",
				Type:      "vcloud",
				Username:  "test",
				Password:  "test",
				VCloudURL: "test",
			}

			data, _ := json.Marshal(mockDC)

			Convey("When I do a post to /datacenters/", func() {
				params := make(map[string]string)
				params["datacenter"] = "test"
				Convey("And I am logged in as an admin", func() {
					resp, err := doRequest("POST", "/datacenters/", params, data, controllers.CreateDatacenterHandler, nil)

					Convey("Then a datacenter should be created", func() {
						var d models.Datacenter
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, 3)
						So(d.Name, ShouldEqual, "new-test")
					})
				})

				SkipConvey("And the datacenter group matches the authenticated users group", func() {
					ft := generateTestToken(1, "test", false)
					resp, err := doRequest("POST", "/datacenters/", params, data, controllers.CreateDatacenterHandler, ft)

					Convey("It should create the datacenter and return the correct set of data", func() {
						var d models.Datacenter
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, 3)
						So(d.Name, ShouldEqual, "new-test")
					})
				})
			})
		})
	})

	Convey("Scenario: deleting a datacenter", t, func() {
		Convey("Given a datacenter exists on the store", func() {
			deleteDatacenterSubscriber()
			getDatacenterSubscriber(2)
			findServiceSubscriber()

			Convey("When I call DELETE /datacenters/:datacenter", func() {
				ft := generateTestToken(1, "test", false)

				params := make(map[string]string)
				params["datacenter"] = "1"
				_, err := doRequest("DELETE", "/datacenters/:datacenter", params, nil, controllers.DeleteDatacenterHandler, ft)

				Convey("It should delete the datacenter and return ok", func() {
					So(err.Error(), ShouldEqual, "code=400, message=Existing services are referring to this datacenter.")
				})
			})

		})

	})
}
