/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/ernestio/api-gateway/config"
	"github.com/ernestio/api-gateway/controllers/datacenters"
	"github.com/ernestio/api-gateway/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDatacenters(t *testing.T) {
	testsSetup()
	config.Setup()
	au := models.User{ID: 1, Username: "test", Password: "test1234"}
	Convey("Scenario: getting a list of datacenters", t, func() {
		Convey("Given datacenters exist on the store", func() {
			findDatacenterSubscriber()
			Convey("When I call /datacenters/", func() {
				st, resp := datacenters.List(au)
				Convey("Then I should have a response with existing datacenters", func() {
					var d []models.Datacenter
					err := json.Unmarshal(resp, &d)
					So(st, ShouldEqual, 200)
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
			getDatacenterSubscriber(1)

			Convey("And I call /datacenter/:datacenter on the api", func() {
				st, resp := datacenters.Get(au, "1")

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should get the existing datacenter", func() {
						var d models.Datacenter
						err := json.Unmarshal(resp, &d)

						So(st, ShouldEqual, 200)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, 1)
						So(d.Name, ShouldEqual, "test")
					})
				})

				Convey("When the datacenter group matches the authenticated users group", func() {
					Convey("Then I should get the existing datacenter", func() {
						var d models.Datacenter
						err := json.Unmarshal(resp, &d)
						So(st, ShouldEqual, 200)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, 1)
						So(d.Name, ShouldEqual, "test")
					})
				})

				SkipConvey("When the datacenter group does not match the authenticated users group", func() {
					st, _ := datacenters.Get(au, "2")
					Convey("Then I should get a 404 error as it doesn't exist", func() {
						So(st, ShouldEqual, 404)
					})
				})
			})
		})
	})

	Convey("Scenario: creating a datacenter", t, func() {
		Convey("Given the datacenter does not exist on the store ", func() {
			getNotFoundDatacenterSubscriber(1)
			createDatacenterSubscriber()

			mockDC := models.Datacenter{
				Name:      "new_test",
				Type:      "vcloud",
				Username:  "test",
				Password:  "test",
				VCloudURL: "test",
			}

			data, _ := json.Marshal(mockDC)

			Convey("When I do a post to /datacenters/", func() {
				Convey("And I am logged in as an admin", func() {
					st, resp := datacenters.Create(au, data)
					Convey("Then a datacenter should be created", func() {
						var d models.Datacenter
						err := json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(st, ShouldEqual, 200)
						So(d.ID, ShouldEqual, 3)
						So(d.Name, ShouldEqual, "new_test")
					})
				})

				Convey("And the datacenter group matches the authenticated users group", func() {
					ft := models.User{ID: 1, Username: "test", Admin: false}
					st, resp := datacenters.Create(ft, data)
					Convey("It should create the datacenter and return the correct set of data", func() {
						var d models.Datacenter
						err := json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(st, ShouldEqual, 200)
						So(d.ID, ShouldEqual, 3)
						So(d.Name, ShouldEqual, "new_test")
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
				res := `[{"resource_id":"1","role":"owner"}]`
				foundSubscriber("authorization.find", res, 1)
				ft := models.User{ID: 1, Username: "test", Admin: false}
				st, resp := datacenters.Delete(ft, "1")
				Convey("It should delete the datacenter and return ok", func() {
					So(st, ShouldEqual, 400)
					So(string(resp), ShouldEqual, "Existing environments are referring to this project.")
				})
			})
		})

	})
}
