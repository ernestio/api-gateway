/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/ernestio/api-gateway/config"
	"github.com/ernestio/api-gateway/controllers/users"
	"github.com/ernestio/api-gateway/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUsers(t *testing.T) {
	var err error
	testsSetup()
	config.Setup()
	au := models.User{ID: 1, GroupID: 1, Username: "test", Password: "test1234"}
	other := models.User{ID: 3, GroupID: 2, Username: "other", Password: "test1234"}
	admin := models.User{ID: 2, Username: "admin", Admin: true}

	Convey("Scenario: getting a list of users", t, func() {
		getGroupSubscriber(3)
		findUserSubscriber()
		Convey("When calling /users/ on the api", func() {
			Convey("And I'm authenticated as an admin user", func() {
				st, resp := users.List(admin)
				Convey("It should show all users", func() {
					var u []models.User
					err = json.Unmarshal(resp, &u)
					So(err, ShouldBeNil)
					So(st, ShouldEqual, 200)
					So(len(u), ShouldEqual, 2)
					So(u[0].ID, ShouldEqual, 1)
					So(u[0].Username, ShouldEqual, "test")
				})
			})
			Convey("And I'm authenticated as a non-admin user", func() {
				st, resp := users.List(au)
				Convey("It should return only the users in the same group", func() {
					var u []models.User
					err = json.Unmarshal(resp, &u)
					So(err, ShouldBeNil)
					So(st, ShouldEqual, 200)
					So(len(u), ShouldEqual, 1)
					So(u[0].ID, ShouldEqual, 1)
					So(u[0].Username, ShouldEqual, "test")
					So(u[0].Password, ShouldEqual, "")
					So(u[0].Salt, ShouldEqual, "")
				})
			})
		})
	})

	Convey("Scenario: getting a single user", t, func() {
		Convey("Given a user exists on the store", func() {
			getUserSubscriber(1)
			Convey("When I call /users/:user on the api", func() {
				Convey("And I'm authenticated as an admin user", func() {
					st, resp := users.Get(admin, "1")
					Convey("It should return the correct set of data", func() {
						var u models.User
						err = json.Unmarshal(resp, &u)
						So(err, ShouldBeNil)
						So(st, ShouldEqual, 200)
						So(u.ID, ShouldEqual, 1)
						So(u.Username, ShouldEqual, "test")
						So(u.Password, ShouldEqual, "")
						So(u.Salt, ShouldEqual, "")
					})
				})
				Convey("And the user is in the same group as a normal user", func() {
					st, resp := users.Get(au, "1")
					Convey("It should return the correct set of data", func() {
						var u models.User
						err = json.Unmarshal(resp, &u)
						So(err, ShouldBeNil)
						So(st, ShouldEqual, 200)
						So(u.ID, ShouldEqual, 1)
						So(u.Username, ShouldEqual, "test")
						So(u.Password, ShouldEqual, "")
						So(u.Salt, ShouldEqual, "")
					})
				})
				Convey("And the user is not in the same group as a normal user", func() {
					st, _ := users.Get(other, "1")
					Convey("It should return a 404", func() {
						So(st, ShouldEqual, 404)
					})
				})
			})
		})

		Convey("Given a user doesn't exist", func() {
			getUserSubscriber(1)
			Convey("When calling /users/:user on the api", func() {
				st, _ := users.Get(au, "99")
				Convey("It should return a 404", func() {
					So(st, ShouldEqual, 404)
				})
			})
		})
	})

	Convey("Scenario: creating a user", t, func() {
		setUserSubscriber()
		getGroupSubscriber(1)
		getUserSubscriber(1)
		Convey("Given no existing users on the store", func() {
			data := []byte(`{"group_id": 1, "username": "new-test", "password": "test1234"}`)
			Convey("When I create a user by calling /users/ on the api", func() {
				Convey("And I'm authenticated as an admin user", func() {
					Convey("With a valid payload", func() {
						st, resp := users.Create(admin, data)
						Convey("It should create the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(err, ShouldBeNil)
							So(st, ShouldEqual, 200)
							So(u.ID, ShouldEqual, 3)
							So(u.Username, ShouldEqual, "new-test")
							So(u.Password, ShouldEqual, "")
							So(u.Salt, ShouldEqual, "")
						})
					})
					Convey("With an invalid payload", func() {
						invalidData := []byte(`{"group_id": 1, "username": "fail"}`)
						st, _ := users.Create(admin, invalidData)
						Convey("It should error with 400 bad request", func() {
							So(st, ShouldEqual, 400)
						})
					})
					Convey("With a password less than the minimum length", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new-test", "password": "test"}`)
						st, resp := users.Create(admin, invalidData)

						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldEqual, `Minimum password length is 8 characters`)
						})
					})
					Convey("With a username using invalid characters", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new^test", "password": "test1234"}`)
						st, resp := users.Create(admin, invalidData)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Username can only contain the following characters: a-z 0-9 @._-")
						})
					})
					Convey("With a password using invalid characters", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new-test", "password": "test^1234"}`)
						st, resp := users.Create(admin, invalidData)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Password can only contain the following characters: a-z 0-9 @._-")
						})
					})
					Convey("With no username", func() {
						invalidData := []byte(`{"group_id": 1, "username": "", "password": "test1234"}`)
						st, resp := users.Create(admin, invalidData)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Username cannot be empty")
						})
					})
					Convey("With no password", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new-test", "password": ""}`)
						st, resp := users.Create(admin, invalidData)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Password cannot be empty")
						})
					})
				})
				Convey("And I'm authenticated as a non-admin user", func() {
					st, _ := users.Create(au, data)
					Convey("It should return with 403 unauthorized", func() {
						So(st, ShouldEqual, 403)
					})
				})
			})
		})

		Convey("Given an existing user on the store", func() {
			existingData := []byte(`{"group_id": 1, "username": "test", "password": "test1234"}`)
			Convey("When I create a user by calling /users/ on the api", func() {
				Convey("And the user already exists", func() {
					st, _ := users.Create(admin, existingData)
					Convey("It should return with 409", func() {
						So(st, ShouldEqual, 409)
					})
				})
			})
		})
	})

	Convey("Scenario: updating a user", t, func() {
		setUserSubscriber()
		getUserSubscriber(1)
		Convey("Given existing users on the store", func() {
			data := []byte(`{"id": 1, "group_id": 1, "username": "test", "password": "new-password"}`)
			Convey("When I update a user by calling /users/ on the api", func() {
				Convey("And I'm authenticated as an admin user", func() {
					Convey("With a valid payload", func() {
						st, resp := users.Update(au, "1", data)
						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(err, ShouldBeNil)
							So(st, ShouldEqual, 200)
							So(u.ID, ShouldEqual, 1)
							So(u.GroupID, ShouldEqual, 1)
							So(u.Username, ShouldEqual, "test")
							So(u.Password, ShouldEqual, "")
							So(u.Salt, ShouldEqual, "")
						})
					})
					Convey("With an invalid payload", func() {
						invalidData := []byte(`{"id": 1, "group_id": 1, "password": "new-password"}`)
						st, _ := users.Create(admin, invalidData)
						Convey("It should update the user and return the correct set of data", func() {
							So(st, ShouldEqual, 400)
						})
					})
					Convey("With a password less than the minimum length", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new-test", "password": "test"}`)
						st, resp := users.Create(admin, invalidData)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldEqual, `Minimum password length is 8 characters`)
						})
					})
					Convey("With a username using invalid characters", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new^test", "password": "test1234"}`)
						st, resp := users.Create(admin, invalidData)

						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Username can only contain the following characters: a-z 0-9 @._-")
						})
					})
					Convey("With a password using invalid characters", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new-test", "password": "test^1234"}`)
						st, resp := users.Create(admin, invalidData)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Password can only contain the following characters: a-z 0-9 @._-")
						})
					})
					Convey("With no username", func() {
						invalidData := []byte(`{"group_id": 1, "username": "", "password": "test1234"}`)
						st, resp := users.Create(admin, invalidData)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Username cannot be empty")
						})
					})
					Convey("With no password", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new-test", "password": ""}`)
						st, resp := users.Create(admin, invalidData)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Password cannot be empty")
						})
					})
					SkipConvey("With an payload id that does not match the user's id", func() {
						//TODO: Finish this.
					})
				})

				Convey("And I'm authenticated as the user being updated", func() {
					Convey("With a valid payload", func() {
						st, resp := users.Update(au, "1", data)
						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(err, ShouldBeNil)
							So(st, ShouldEqual, 200)
							So(u.ID, ShouldEqual, 1)
							So(u.GroupID, ShouldEqual, 1)
							So(u.Username, ShouldEqual, "test")
							So(u.Password, ShouldEqual, "")
							So(u.Salt, ShouldEqual, "")
						})
					})
					Convey("With a group id that does not match the exisiting users id", func() {
						invalidData := []byte(`{"id": 1, "group_id": 2, "username": "test", "password": "new-password"}`)
						st, _ := users.Update(au, "1", invalidData)
						Convey("It should update the user and return the correct set of data", func() {
							So(st, ShouldEqual, 403)
						})
					})
				})

				Convey("And I'm not authenticated as the user being updated", func() {
					st, _ := users.Update(other, "1", data)
					Convey("It should return with 403 unauthorized", func() {
						So(st, ShouldEqual, 403)
					})
				})
			})
		})

		Convey("Given no existing users on the store", func() {
			data := []byte(`{"id": 99, "group_id": 1, "username": "fake-user", "password": "test1234"}`)

			Convey("And I update a user by calling /users/ on the api", func() {
				st, _ := users.Update(admin, "99", data)

				Convey("It should error with 404 doesn't exist", func() {
					So(st, ShouldEqual, 404)
				})
			})
		})

	})

	Convey("Scenario: deleting a user", t, func() {
		deleteUserSubscriber()

		Convey("Given existing users on the store", func() {
			Convey("When I delete a user by calling /users/:user on the api", func() {
				Convey("And I am logged in as an admin", func() {
					st, _ := users.Delete(admin, "1")

					Convey("It should delete the user and return a 200 ok", func() {
						So(st, ShouldEqual, 200)
					})
				})
				Convey("And I am logged in as a non-admin", func() {
					st, _ := users.Delete(au, "1")

					Convey("It should return a 403 not authorized", func() {
						So(st, ShouldEqual, 403)
					})
				})
			})
		})
		Convey("Given no users on the store", func() {
			Convey("When I delete a user by calling /users/:user on the api", func() {
				st, _ := users.Delete(admin, "99")

				Convey("It should return a 404 ok", func() {
					So(st, ShouldEqual, 404)
				})
			})
		})
	})
}
