/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGroups(t *testing.T) {
	testsSetup()
	setup()

	Convey("Scenario: getting a list of groups", t, func() {
		Convey("Given groups exist on the store", func() {
			findGroupSubscriber()
			Convey("When I call /groups/", func() {
				resp, err := doRequest("GET", "/groups/", nil, nil, getGroupsHandler, nil)
				Convey("Then I should have a response existing groups", func() {
					var g []Group
					So(err, ShouldBeNil)

					err = json.Unmarshal(resp, &g)

					So(err, ShouldBeNil)
					So(len(g), ShouldEqual, 2)
					So(g[0].ID, ShouldEqual, 1)
					So(g[0].Name, ShouldEqual, "test")
				})
			})

			SkipConvey("Given no groups on the store", func() {
			})
		})
	})

	Convey("Scenario: getting a single group", t, func() {
		Convey("Given the group exist on the store", func() {
			getGroupSubscriber()
			Convey("And I call /groups/:group on the api", func() {
				params := make(map[string]string)
				params["group"] = "1"
				resp, err := doRequest("GET", "/groups/:group", params, nil, getGroupHandler, nil)

				Convey("When I'm authenticated as admin user", func() {
					Convey("Then I should get the existing group", func() {
						var g Group

						So(err, ShouldBeNil)

						err = json.Unmarshal(resp, &g)

						So(err, ShouldBeNil)
						So(g.ID, ShouldEqual, 1)
						So(g.Name, ShouldEqual, "test")
					})
				})
			})
		})
	})

	Convey("Scenario: create a group", t, func() {
		Convey("Given a group exists on the store ", func() {
			createGroupSubscriber()
			getGroupSubscriber()

			mockG := Group{
				ID:   1,
				Name: "new-test",
			}

			data, _ := json.Marshal(mockG)

			Convey("When I do a post to /groups/", func() {
				params := make(map[string]string)
				params["group"] = "test"
				resp, err := doRequest("POST", "/groups/", params, data, createGroupHandler, nil)
				Convey("Then a group hould be created", func() {
					var g Group
					So(err, ShouldBeNil)
					err = json.Unmarshal(resp, &g)
					So(err, ShouldBeNil)
					So(g.ID, ShouldEqual, 3)
					So(g.Name, ShouldEqual, "new-test")
				})
			})
		})
	})

	Convey("Scenario: deleting a group", t, func() {
		Convey("Given a group exists on the store", func() {
			deleteGroupSubscriber()
			getGroupSubscriber()

			Convey("When I call DELETE /groups/:group", func() {
				Convey("And I am logged in as an admin", func() {
					params := make(map[string]string)
					params["group"] = "test"
					_, err := doRequest("DELETE", "/groups/:group", params, nil, deleteGroupHandler, nil)

					SkipConvey("It should delete the group and return ok", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		})
	})
}
