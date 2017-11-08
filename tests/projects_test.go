/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/ernestio/api-gateway/config"
	"github.com/ernestio/api-gateway/controllers/projects"
	"github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestListProjects(t *testing.T) {
	testsSetup()
	config.Setup()
	au := models.User{ID: 1, Username: "test", Password: &pw1}
	Convey("Scenario: getting a list of datacenters", t, func() {
		Convey("Given datacenters exist on the store", func() {
			findDatacenterSubscriber()
			Convey("When I call /datacenters/", func() {
				st, resp := projects.List(au)
				Convey("Then I should have a response with existing datacenters", func() {
					var d []models.Project
					err := json.Unmarshal(resp, &d)
					So(st, ShouldEqual, 200)
					So(err, ShouldBeNil)
					So(len(d), ShouldEqual, 2)
					So(d[0].ID, ShouldEqual, 1)
					So(d[0].Name, ShouldEqual, "test")
				})
			})
		})
	})
}

func TestGetProject(t *testing.T) {
	testsSetup()
	config.Setup()
	au := models.User{ID: 1, Username: "test", Password: &pw1}
	Convey("Scenario: getting a single datacenters", t, func() {
		Convey("Given the datacenter exists on the store", func() {
			getDatacenterSubscriber(1)

			Convey("And I call /datacenter/:datacenter on the api", func() {
				st, resp := projects.Get(au, "1")

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should get the existing datacenter", func() {
						var d models.Project
						err := json.Unmarshal(resp, &d)

						So(st, ShouldEqual, 200)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, 1)
					})
				})
				Convey("When the datacenter group matches the authenticated users group", func() {
					Convey("Then I should get the existing datacenter", func() {
						var d models.Project
						err := json.Unmarshal(resp, &d)
						So(st, ShouldEqual, 200)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, 1)
						So(d.Name, ShouldEqual, "test")
					})
				})
			})
		})
	})
}

func TestCreateProject(t *testing.T) {
	testsSetup()
	config.Setup()
	au := models.User{ID: 1, Username: "test", Password: &pw1}

	Convey("Scenario: creating a datacenter", t, func() {
		Convey("Given the datacenter does not exist on the store ", func() {
			getNotFoundDatacenterSubscriber(1)
			createDatacenterSubscriber()

			mockDC := models.Project{
				Name: "new_test",
				Type: "vcloud",
				Credentials: map[string]interface{}{
					"username":   "test",
					"password":   "test",
					"vcloud_url": "test",
				},
			}

			data, _ := json.Marshal(mockDC)

			Convey("When I do a post to /datacenters/", func() {
				Convey("And I am logged in as an admin", func() {
					st, resp := projects.Create(au, data)
					Convey("Then a datacenter should be created", func() {
						var d models.Project
						err := json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(st, ShouldEqual, 200)
						So(d.ID, ShouldEqual, 3)
						So(d.Name, ShouldEqual, "new_test")
					})
				})

				Convey("And the datacenter group matches the authenticated users group", func() {
					ft := models.User{ID: 1, Username: "test", Admin: helpers.Bool(false)}
					st, resp := projects.Create(ft, data)
					Convey("It should create the datacenter and return the correct set of data", func() {
						var d models.Project
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
}

func TestDeleteProject(t *testing.T) {
	testsSetup()
	config.Setup()
	Convey("Scenario: deleting a datacenter", t, func() {
		Convey("Given a datacenter exists on the store", func() {
			getDatacenterSolo(1)
			findServiceSolo()

			Convey("When I call DELETE /datacenters/:datacenter", func() {
				foundSubscriber("authorization.find", `[{"resource_id":"1","role":"owner"}]`, 1)
				ft := models.User{ID: 1, Username: "test", Admin: helpers.Bool(false)}
				st, resp := projects.Delete(ft, "1")
				Convey("It should delete the datacenter and return ok", func() {
					So(st, ShouldEqual, 400)
					So(string(resp), ShouldEqual, "Existing environments are referring to this project.")
				})
			})
		})
	})
}
