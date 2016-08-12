/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUsers(t *testing.T) {
	testsSetup()
	setup()

	Convey("Scenario: getting a list of users", t, func() {
		findUserSubscriber()
		Convey("When calling /users/ on the api", func() {
			Convey("And I'm authenticated as an admin user", func() {
				params := make(map[string]string)
				ft := generateTestToken(1, "admin", true)
				resp, err := doRequest("GET", "/users/", params, nil, getUsersHandler, ft)
				Convey("It should show all users", func() {
					var u []User

					So(err, ShouldBeNil)

					err = json.Unmarshal(resp, &u)

					So(err, ShouldBeNil)
					So(len(u), ShouldEqual, 2)
					So(u[0].ID, ShouldEqual, 1)
					So(u[0].Username, ShouldEqual, "test")
				})
			})
			Convey("And I'm authenticated as a non-admin user", func() {
				params := make(map[string]string)
				ft := generateTestToken(1, "test", false)
				resp, err := doRequest("GET", "/users/", params, nil, getUsersHandler, ft)

				Convey("It should return only the users in the same group", func() {
					var u []User

					So(err, ShouldBeNil)

					err = json.Unmarshal(resp, &u)

					So(err, ShouldBeNil)
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
					params := make(map[string]string)
					params["user"] = "1"
					ft := generateTestToken(1, "admin", true)
					resp, err := doRequest("GET", "/users/:user", params, nil, getUserHandler, ft)

					Convey("It should return the correct set of data", func() {
						var u User

						So(err, ShouldBeNil)

						err = json.Unmarshal(resp, &u)

						So(err, ShouldBeNil)
						So(u.ID, ShouldEqual, 1)
						So(u.Username, ShouldEqual, "test")
						So(u.Password, ShouldEqual, "")
						So(u.Salt, ShouldEqual, "")
					})
				})
				Convey("And the user is in the same group as a normal user", func() {
					params := make(map[string]string)
					params["user"] = "1"
					ft := generateTestToken(1, "test", false)
					resp, err := doRequest("GET", "/users/:user", params, nil, getUserHandler, ft)

					Convey("It should return the correct set of data", func() {
						var u User

						So(err, ShouldBeNil)

						err = json.Unmarshal(resp, &u)

						So(err, ShouldBeNil)
						So(u.ID, ShouldEqual, 1)
						So(u.Username, ShouldEqual, "test")
						So(u.Password, ShouldEqual, "")
						So(u.Salt, ShouldEqual, "")
					})
				})
				Convey("And the user is not in the same group as a normal user", func() {
					params := make(map[string]string)
					params["user"] = "1"
					ft := generateTestToken(2, "test2", false)
					resp, err := doRequest("GET", "/users/:user", params, nil, getUserHandler, ft)

					Convey("It should return a 404", func() {
						So(err, ShouldNotBeNil)
						So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
						So(len(resp), ShouldEqual, 0)
					})
				})
			})
		})

		Convey("Given a user doesn't exist", func() {
			getUserSubscriber(1)
			Convey("When calling /users/:user on the api", func() {
				params := make(map[string]string)
				params["user"] = "99"
				ft := generateTestToken(2, "test2", false)
				resp, err := doRequest("GET", "/users/:user", params, nil, getUserHandler, ft)

				Convey("It should return a 404", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
					So(len(resp), ShouldEqual, 0)
				})
			})
		})
	})

	Convey("Scenario: creating a user", t, func() {
		setUserSubscriber()
		getGroupSubscriber()
		getUserSubscriber(1)
		Convey("Given no existing users on the store", func() {
			data := []byte(`{"group_id": 1, "username": "new-test", "password": "test"}`)

			Convey("When I create a user by calling /users/ on the api", func() {
				Convey("And I'm authenticated as an admin user", func() {
					Convey("With a valid payload", func() {
						ft := generateTestToken(1, "admin", true)
						resp, err := doRequest("POST", "/users/", nil, data, createUserHandler, ft)

						Convey("It should create the user and return the correct set of data", func() {
							var u User

							So(err, ShouldBeNil)

							err = json.Unmarshal(resp, &u)

							So(err, ShouldBeNil)
							So(u.ID, ShouldEqual, 3)
							So(u.Username, ShouldEqual, "new-test")
							So(u.Password, ShouldEqual, "")
							So(u.Salt, ShouldEqual, "")
						})
					})
					Convey("With an invalid payload", func() {
						invalidData := []byte(`{"group_id": 1, "username": "fail"}`)
						ft := generateTestToken(1, "admin", true)
						_, err := doRequest("POST", "/users/", nil, invalidData, createUserHandler, ft)

						Convey("It should error with 400 bad request", func() {
							So(err, ShouldNotBeNil)
							So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
						})
					})
				})
				Convey("And I'm authenticated as a non-admin user", func() {
					ft := generateTestToken(1, "test2", false)
					_, err := doRequest("POST", "/users/", nil, data, createUserHandler, ft)
					Convey("It should return with 403 unauthorized", func() {
						So(err, ShouldNotBeNil)
						So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
					})
				})
			})

		})

		Convey("Given an existing user on the store", func() {
			existingData := []byte(`{"group_id": 1, "username": "test", "password": "test"}`)
			Convey("When I create a user by calling /users/ on the api", func() {
				Convey("And the user already exists", func() {
					ft := generateTestToken(1, "admin", true)
					_, err := doRequest("POST", "/users/", nil, existingData, createUserHandler, ft)

					Convey("It should return with 303 see other", func() {
						So(err, ShouldNotBeNil)
						So(err.(*echo.HTTPError).Code, ShouldEqual, 303)
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
					params := make(map[string]string)
					params["user"] = "1"
					ft := generateTestToken(1, "admin", true)
					Convey("With a valid payload", func() {
						resp, err := doRequest("PUT", "/users/:user", params, data, updateUserHandler, ft)
						Convey("It should update the user and return the correct set of data", func() {
							var u User

							So(err, ShouldBeNil)

							err = json.Unmarshal(resp, &u)

							So(err, ShouldBeNil)
							So(u.ID, ShouldEqual, 1)
							So(u.GroupID, ShouldEqual, 1)
							So(u.Username, ShouldEqual, "test")
							So(u.Password, ShouldEqual, "")
							So(u.Salt, ShouldEqual, "")
						})
					})
					Convey("With an invalid payload", func() {
						invalidData := []byte(`{"id": 1, "group_id": 1, "password": "new-password"}`)
						_, err := doRequest("PUT", "/users/:user", params, invalidData, updateUserHandler, ft)
						Convey("It should update the user and return the correct set of data", func() {
							So(err, ShouldNotBeNil)
							So(err.(*echo.HTTPError).Code, ShouldEqual, 400)
						})
					})
					SkipConvey("With an payload id that does not match the user's id", func() {
						//TODO: Finish this.
					})
				})

				Convey("And I'm authenticated as the user being updated", func() {
					Convey("With a valid payload", func() {
						params := make(map[string]string)
						params["user"] = "1"
						ft := generateTestToken(1, "test", false)
						resp, err := doRequest("PUT", "/users/:user", params, data, updateUserHandler, ft)
						Convey("It should update the user and return the correct set of data", func() {
							var u User

							So(err, ShouldBeNil)

							err = json.Unmarshal(resp, &u)

							So(err, ShouldBeNil)
							So(u.ID, ShouldEqual, 1)
							So(u.GroupID, ShouldEqual, 1)
							So(u.Username, ShouldEqual, "test")
							So(u.Password, ShouldEqual, "")
							So(u.Salt, ShouldEqual, "")
						})
					})
					Convey("With a group id that does not match the exisiting users id", func() {
						invalidData := []byte(`{"id": 1, "group_id": 2, "username": "test", "password": "new-password"}`)
						params := make(map[string]string)
						params["user"] = "1"
						ft := generateTestToken(1, "test", false)
						_, err := doRequest("PUT", "/users/:user", params, invalidData, updateUserHandler, ft)
						Convey("It should update the user and return the correct set of data", func() {
							So(err, ShouldNotBeNil)
							So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
						})
					})
				})

				Convey("And I'm not authenticated as the user being updated", func() {
					ft := generateTestToken(1, "test2", false)
					params := make(map[string]string)
					params["user"] = "2"
					_, err := doRequest("PUT", "/users/:user", params, data, updateUserHandler, ft)

					Convey("It should return with 403 unauthorized", func() {
						So(err, ShouldNotBeNil)
						So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
					})
				})
			})
		})

		Convey("Given no existing users on the store", func() {
			data := []byte(`{"id": 99, "group_id": 1, "username": "fake-user", "password": "test"}`)

			Convey("And I update a user by calling /users/ on the api", func() {
				ft := generateTestToken(1, "admin", true)
				params := make(map[string]string)
				params["user"] = "99"
				_, err := doRequest("PUT", "/users/:user", params, data, updateUserHandler, ft)

				Convey("It should error with 404 doesn't exist", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
				})
			})
		})

	})

	Convey("Scenario: deleting a user", t, func() {
		deleteUserSubscriber()

		Convey("Given existing users on the store", func() {
			Convey("When I delete a user by calling /users/:user on the api", func() {
				Convey("And I am logged in as an admin", func() {
					ft := generateTestToken(1, "admin", true)
					params := make(map[string]string)
					params["user"] = "1"
					_, err := doRequest("DELETE", "/users/:user", params, nil, deleteUserHandler, ft)

					Convey("It should delete the user and return a 200 ok", func() {
						So(err, ShouldBeNil)
					})
				})
				Convey("And I am logged in as a non-admin", func() {
					ft := generateTestToken(1, "test", false)
					params := make(map[string]string)
					params["user"] = "1"
					_, err := doRequest("DELETE", "/users/:user", params, nil, deleteUserHandler, ft)

					Convey("It should return a 403 not authorized", func() {
						So(err, ShouldNotBeNil)
						So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
					})
				})
			})
		})
		Convey("Given no users on the store", func() {
			Convey("When I delete a user by calling /users/:user on the api", func() {
				ft := generateTestToken(1, "admin", true)
				params := make(map[string]string)
				params["user"] = "99"
				_, err := doRequest("DELETE", "/users/:user", params, nil, deleteUserHandler, ft)

				Convey("It should return a 404 ok", func() {
					So(err, ShouldNotBeNil)
					So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
				})
			})
		})
	})
}
