/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/ernestio/api-gateway/controllers/users"
	"github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetUsers(t *testing.T) {
	var err error
	testsSetup()
	au := models.User{ID: 1, Username: "test", Password: &pw1}
	admin := models.User{ID: 2, Username: "admin", Admin: helpers.Bool(true)}

	Convey("Scenario: getting a list of users", t, func() {
		findUserSubscriber()
		Convey("When calling /users/ on the api", func() {
			Convey("And I'm authenticated as an admin user", func() {
				foundSubscriber("user.get", `{"id":"1","name":"fake/test","datacenter_id":1}`, 1)
				st, resp := users.List(admin)
				Convey("It should show all users", func() {
					var u []models.User
					err = json.Unmarshal(resp, &u)
					So(err, ShouldBeNil)
					So(st, ShouldEqual, 200)
					So(len(u), ShouldEqual, 2)
					So(u[0].ID, ShouldEqual, 1)
					So(u[0].Username, ShouldEqual, "test")
					So(u[0].Password, ShouldNotBeNil)
					So(u[0].ProjectMemberships, ShouldBeEmpty)
					So(u[0].EnvMemberships, ShouldBeEmpty)
					So(u[0].Type, ShouldBeBlank)
				})
			})
			Convey("And I'm authenticated as a non-admin user", func() {
				st, resp := users.List(au)
				Convey("It should show all users with limited information", func() {
					var u []models.User
					err = json.Unmarshal(resp, &u)
					So(err, ShouldBeNil)
					So(st, ShouldEqual, 200)
					So(len(u), ShouldEqual, 2)
					So(u[0].ID, ShouldEqual, 1)
					So(u[0].Username, ShouldEqual, "test")
					So(u[0].Password, ShouldBeNil)
					So(u[0].ProjectMemberships, ShouldBeNil)
					So(u[0].EnvMemberships, ShouldBeNil)
					So(u[0].Admin, ShouldBeNil)
					So(u[0].Type, ShouldBeBlank)
					So(u[0].Disabled, ShouldBeNil)
				})
			})
		})
	})
}

func TestGetUser(t *testing.T) {
	var err error
	testsSetup()
	au := models.User{ID: 1, Username: "test", Password: &pw1}
	other := models.User{ID: 3, Username: "other", Password: &pw1}
	admin := models.User{ID: 2, Username: "admin", Admin: helpers.Bool(true)}

	Convey("Scenario: getting a single user", t, func() {
		Convey("Given a user exists on the store", func() {
			getUserSubscriber(1)
			foundSubscriber("authorization.find", `[{"role":"owner"}]`, 2)
			Convey("When I call /users/:user on the api", func() {
				Convey("And I'm authenticated as an admin user", func() {
					st, resp := users.Get(admin, "test")
					Convey("It should return the correct set of data", func() {
						var u models.User
						err = json.Unmarshal(resp, &u)
						So(st, ShouldEqual, 200)
						So(err, ShouldBeNil)
						So(u.ID, ShouldEqual, 1)
						So(u.Username, ShouldEqual, "test")
						So(*u.Password, ShouldBeBlank)
						So(u.Salt, ShouldBeBlank)
						So(u.MFASecret, ShouldBeBlank)
					})
				})
				Convey("And the user is the same as registered  user", func() {
					st, resp := users.Get(au, "test")
					Convey("It should return the correct set of data", func() {
						var u models.User
						err = json.Unmarshal(resp, &u)
						So(st, ShouldEqual, 200)
						So(err, ShouldBeNil)
						So(u.ID, ShouldEqual, 1)
						So(u.Username, ShouldEqual, "test")
						So(u.Password, ShouldBeNil)
						So(u.Salt, ShouldEqual, "")
						So(u.MFASecret, ShouldBeBlank)
					})
				})
				Convey("And the user is not the same as a registered user", func() {
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
}

func TestCreateUser(t *testing.T) {
	var err error
	testsSetup()
	admin := models.User{ID: 2, Username: "admin", Admin: helpers.Bool(true)}

	Convey("Scenario: creating a user", t, func() {
		setUserSubscriber()
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
							So(*u.Password, ShouldBeBlank)
							So(u.Salt, ShouldBeBlank)
							So(u.MFASecret, ShouldBeBlank)
						})
					})
					Convey("With a password less than the minimum length", func() {
						invalidData := []byte(`{"group_id": 1, "username": "new-test", "password": "test"}`)
						st, resp := users.Create(admin, invalidData)

						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldEqual, `{"message":"Minimum password length is 8 characters"}`)
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
}

func TestUpdateUser(t *testing.T) {
	var err error
	testsSetup()
	au := models.User{ID: 1, Username: "test", Password: &pw1}
	other := models.User{ID: 3, Username: "other", Password: &pw1}
	admin := models.User{ID: 2, Username: "admin", Admin: helpers.Bool(true)}
	data := []byte(`{"id": 1, "username": "test", "password": "new-password"}`)

	Convey("Scenario: updating a user", t, func() {
		setUserSubscriber()
		Convey("Given existing users on the store", func() {
			Convey("When I update a user by calling /users/ on the api", func() {
				Convey("And I'm authenticated as an admin user", func() {
					Convey("With a valid payload", func() {
						var err error
						getUserSubscriber(1)
						st, resp := users.Update(admin, "test", data)
						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(err, ShouldBeNil)
							So(st, ShouldEqual, 200)
							So(u.ID, ShouldEqual, 1)
							So(u.Username, ShouldEqual, "test")
							So(*u.Password, ShouldBeBlank)
							So(u.Salt, ShouldBeBlank)
							So(u.MFASecret, ShouldBeBlank)
						})
					})
					Convey("With a payload updating the admin field", func() {
						data := []byte(`{"id": 1, "username": "test", "password": "new-password", "admin": true}`)
						getUserSubscriber(1)
						st, resp := users.Update(admin, "test", data)
						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(err, ShouldBeNil)
							So(st, ShouldEqual, 200)
							So(u.ID, ShouldEqual, 1)
							So(u.Username, ShouldEqual, "test")
							So(*u.Password, ShouldBeBlank)
							So(u.Salt, ShouldBeBlank)
							So(*u.Admin, ShouldEqual, true)
							So(u.MFASecret, ShouldBeBlank)
						})
					})
					Convey("With a payload containing no username", func() {
						getUserSubscriber(1)
						data := []byte(`{"id": 1, "username": "test", "password": "new-password"}`)
						st, _ := users.Update(admin, "test", data)
						Convey("It should update the user and return the correct set of data", func() {
							So(st, ShouldEqual, 200)
						})
					})
					Convey("With a payload updating the username", func() {
						data := []byte(`{"id": 1, "username": "newtest", "password": "new-password"}`)
						getUserSubscriber(1)
						st, resp := users.Update(admin, "newtest", data)

						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 404)
							So(string(resp), ShouldContainSubstring, "Specified user not found")
						})
					})
					Convey("With a password using invalid characters", func() {
						getUserSubscriber(1)
						data := []byte(`{"id": 1, "username": "test", "password": "new^password"}`)
						st, resp := users.Update(admin, "test", data)
						Convey("It should return an error message with a 400 repsonse", func() {
							So(st, ShouldEqual, 400)
							So(string(resp), ShouldContainSubstring, "Password can only contain the following characters: a-z 0-9 @._-")
						})
					})
					Convey("With no password", func() {
						getUserSubscriber(1)
						data := []byte(`{"id": 1, "username": "test", "admin": true}`)
						st, resp := users.Update(admin, "test", data)

						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(st, ShouldEqual, 200)
							So(*u.Admin, ShouldEqual, true)
						})
					})
					Convey("With a payload enabling MFA", func() {
						getUserSubscriber(1)
						data := []byte(`{"id": 1, "username": "test", "mfa": true}`)
						st, resp := users.Update(admin, "test", data)

						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(st, ShouldEqual, 200)
							So(u.Username, ShouldEqual, "test")
							So(*u.MFA, ShouldEqual, true)
							So(u.MFASecret, ShouldNotBeBlank)
						})
					})
					Convey("With a payload disabling MFA", func() {
						getUserSubscriber(1)
						data := []byte(`{"id": 1, "username": "test", "mfa": false}`)
						st, resp := users.Update(admin, "test", data)

						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(st, ShouldEqual, 200)
							So(u.Username, ShouldEqual, "test")
							So(*u.MFA, ShouldEqual, false)
							So(u.MFASecret, ShouldBeBlank)
						})
					})
				})
				Convey("And I'm authenticated as the user being updated", func() {
					Convey("With a valid payload", func() {
						getUserSubscriber(1)
						st, resp := users.Update(au, "test", data)
						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(err, ShouldBeNil)
							So(st, ShouldEqual, 200)
							So(u.ID, ShouldEqual, 1)
							So(u.Username, ShouldEqual, "test")
							So(u.Password, ShouldBeNil)
							So(u.Salt, ShouldBeBlank)
							So(u.MFASecret, ShouldBeBlank)
						})
					})
					Convey("With a payload updating the admin field", func() {
						data := []byte(`{"id": 1, "username": "test", "password": "new-password", "admin": true}`)
						getUserSubscriber(1)
						st, resp := users.Update(au, "test", data)
						Convey("It should return an error message with a 403 response", func() {
							So(st, ShouldEqual, 403)
							So(string(resp), ShouldContainSubstring, "You're not allowed to perform this action, please contact your admin")
						})
					})
					Convey("With a payload enabling MFA", func() {
						getUserSubscriber(1)
						data := []byte(`{"id": 1, "username": "test", "mfa": true}`)
						st, resp := users.Update(admin, "test", data)

						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(st, ShouldEqual, 200)
							So(u.Username, ShouldEqual, "test")
							So(*u.MFA, ShouldEqual, true)
							So(u.MFASecret, ShouldNotBeBlank)
						})
					})
					Convey("With a payload disabling MFA", func() {
						getUserSubscriber(1)
						data := []byte(`{"id": 1, "username": "test", "mfa": false}`)
						st, resp := users.Update(admin, "test", data)

						Convey("It should update the user and return the correct set of data", func() {
							var u models.User
							err = json.Unmarshal(resp, &u)
							So(st, ShouldEqual, 200)
							So(u.Username, ShouldEqual, "test")
							So(*u.MFA, ShouldEqual, false)
							So(u.MFASecret, ShouldBeBlank)
						})
					})
				})
				Convey("And I'm not authenticated as the user being updated", func() {
					getUserSubscriber(1)
					st, _ := users.Update(other, "test", data)
					Convey("It should return with 403 unauthorized", func() {
						So(st, ShouldEqual, 403)
					})
				})
			})
		})

		Convey("Given no existing users on the store", func() {
			data := []byte(`{"id": 99, "group_id": 1, "username": "fake-user", "password": "test1234"}`)

			Convey("And I update a user by calling /users/ on the api", func() {
				st, _ := users.Update(admin, "fake-user", data)

				Convey("It should error with 404 doesn't exist", func() {
					So(st, ShouldEqual, 404)
				})
			})
		})

	})
}

func TestDeleteUser(t *testing.T) {
	testsSetup()
	admin := models.User{ID: 2, Username: "admin", Admin: helpers.Bool(true)}

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
