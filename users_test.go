/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mockUsers = []User{
		User{
			ID:       "1",
			GroupID:  "1",
			Username: "test",
			Password: "test",
		},
		User{
			ID:       "2",
			GroupID:  "2",
			Username: "test2",
			Password: "test2",
		},
	}
)

func getUsersSubcriber() {
	n.Subscribe("users.get", func(msg *nats.Msg) {
		data, _ := json.Marshal(mockUsers)
		n.Publish(msg.Reply, data)
	})
}

func getUserSubcriber() {
	n.Subscribe("users.get.1", func(msg *nats.Msg) {
		data, _ := json.Marshal(mockUsers[0])
		n.Publish(msg.Reply, data)
	})
}

func findUserSubcriber() {
	n.Subscribe("users.find", func(msg *nats.Msg) {
		var u User
		json.Unmarshal(msg.Data, &u)

		for _, user := range mockUsers {
			if user.Username == u.Username {
				u = user
				break
			}
		}

		data, _ := json.Marshal(u)
		n.Publish(msg.Reply, data)
	})
}

func createUserSubcriber() {
	n.Subscribe("users.create", func(msg *nats.Msg) {
		var u User

		json.Unmarshal(msg.Data, &u)
		u.ID = "3"
		data, _ := json.Marshal(u)

		n.Publish(msg.Reply, data)
	})
}

func deleteUserSubcriber() {
	n.Subscribe("users.delete.1", func(msg *nats.Msg) {
		n.Publish(msg.Reply, []byte{})
	})
}

func TestUsers(t *testing.T) {
	Convey("Given user handler", t, func() {
		// setup nats connection
		os.Setenv("JWT_SECRET", "test")
		setup()

		Convey("When getting a list of users", func() {
			getUsersSubcriber()

			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
			c.SetPath("/users/")

			Convey("It should return the correct set of data", func() {
				var u []User

				err := getUsersHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &u)

				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, 200)
				So(len(u), ShouldEqual, 2)
				So(u[0].ID, ShouldEqual, "1")
				So(u[0].Username, ShouldEqual, "test")
			})

		})

		Convey("When getting a single user", func() {
			getUserSubcriber()

			e := echo.New()
			req := new(http.Request)
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			c.SetPath("/users/:user")
			c.SetParamNames("user")
			c.SetParamValues("1")

			Convey("It should return the correct set of data", func() {
				var u User

				err := getUserHandler(c)
				So(err, ShouldBeNil)

				resp := rec.Body.Bytes()
				err = json.Unmarshal(resp, &u)

				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, 200)
				So(u.ID, ShouldEqual, "1")
				So(u.Username, ShouldEqual, "test")
			})

		})

		Convey("When creating a user", func() {
			createUserSubcriber()

			Convey("With a valid payload", func() {
				data, _ := json.Marshal(User{GroupID: "1", Username: "new-test", Password: "test"})

				Convey("As an admin user", func() {
					e := echo.New()
					req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(data))
					rec := httptest.NewRecorder()
					c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

					ft := jwt.New(jwt.SigningMethodHS256)
					ft.Claims["username"] = "test"
					ft.Claims["admin"] = true

					c.SetPath("/users/")
					c.Set("user", ft)

					Convey("It should create the user and return the correct set of data", func() {
						var u User

						err := createUserHandler(c)
						So(err, ShouldBeNil)

						resp := rec.Body.Bytes()
						err = json.Unmarshal(resp, &u)

						So(err, ShouldBeNil)
						So(u.ID, ShouldEqual, "3")
						So(u.Username, ShouldEqual, "new-test")
					})

				})

				Convey("As an non-admin user", func() {
					e := echo.New()
					req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(data))
					rec := httptest.NewRecorder()
					c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

					ft := jwt.New(jwt.SigningMethodHS256)
					ft.Claims["username"] = "test"
					ft.Claims["admin"] = false

					c.SetPath("/users/")
					c.Set("user", ft)

					Convey("It should return with 403 unauthorized", func() {
						err := createUserHandler(c).(*echo.HTTPError)
						So(err, ShouldNotBeNil)
						So(err.Code, ShouldEqual, 403)
					})
				})

			})

			Convey("With an invalid payload", func() {
				data := []byte(`{"group_id": 1, "username": "fail"}`)

				e := echo.New()
				req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(data))
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "test"
				ft.Claims["admin"] = true

				c.Set("user", ft)
				c.SetPath("/users/")

				Convey("It should error with 400 bad request", func() {
					err := createUserHandler(c).(*echo.HTTPError)
					So(err, ShouldNotBeNil)
					So(err.Code, ShouldEqual, 400)
				})
			})

		})

		Convey("When updating a user", func() {
			// TODO : Uncomment this
			// updateUserSubcriber()

			Convey("Thats exists", func() {
				Convey("As an admin user", func() {
					data, _ := json.Marshal(User{GroupID: "1", Username: "test2", Password: "test"})

					e := echo.New()
					req, _ := http.NewRequest("POST", "/users/1", bytes.NewReader(data))
					rec := httptest.NewRecorder()
					c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

					ft := jwt.New(jwt.SigningMethodHS256)
					ft.Claims["username"] = "admin"
					ft.Claims["admin"] = true

					c.SetPath("/users/:user")
					c.SetParamNames("user")
					c.SetParamValues("1")
					c.Set("user", ft)

					Convey("It should update the user and return the correct set of data", func() {
						var u User

						err := updateUserHandler(c)
						So(err, ShouldBeNil)

						resp := rec.Body.Bytes()
						err = json.Unmarshal(resp, &u)

						So(err, ShouldBeNil)
						So(u.ID, ShouldEqual, "2")
						So(u.GroupID, ShouldEqual, "1")
						So(u.Username, ShouldEqual, "test2")
					})

				})

				Convey("As an non-admin user", func() {
					Convey("Where a user updates itself", func() {
						data, _ := json.Marshal(User{GroupID: "1", Username: "test", Password: "test2"})

						e := echo.New()
						req, _ := http.NewRequest("POST", "/users/1", bytes.NewReader(data))
						rec := httptest.NewRecorder()
						c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

						ft := jwt.New(jwt.SigningMethodHS256)
						ft.Claims["username"] = "test"
						ft.Claims["admin"] = false

						c.SetPath("/users/:user")
						c.SetParamNames("user")
						c.SetParamValues("1")
						c.Set("user", ft)

						Convey("It should update the user and return the correct set of data", func() {
							var u User

							err := updateUserHandler(c)
							So(err, ShouldBeNil)

							resp := rec.Body.Bytes()
							err = json.Unmarshal(resp, &u)

							So(err, ShouldBeNil)
							So(u.ID, ShouldEqual, "2")
							So(u.GroupID, ShouldEqual, "1")
							So(u.Username, ShouldEqual, "test2")
						})
					})
					Convey("Where a user updates another user", func() {
						data, _ := json.Marshal(User{GroupID: "1", Username: "test2", Password: "test2"})

						e := echo.New()
						req, _ := http.NewRequest("POST", "/users/1", bytes.NewReader(data))
						rec := httptest.NewRecorder()
						c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

						ft := jwt.New(jwt.SigningMethodHS256)
						ft.Claims["username"] = "test"
						ft.Claims["admin"] = false

						c.SetPath("/users/:user")
						c.SetParamNames("user")
						c.SetParamValues("1")
						c.Set("user", ft)

						Convey("It should return with 403 unauthorized", func() {
							err := updateUserHandler(c).(*echo.HTTPError)
							So(err, ShouldNotBeNil)
							So(err.Code, ShouldEqual, 403)
						})
					})
				})

			})

			Convey("That doesn't exist", func() {
				data := []byte(`{"group_id": 1, "username": "fail"}`)

				e := echo.New()
				req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(data))
				rec := httptest.NewRecorder()
				c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

				ft := jwt.New(jwt.SigningMethodHS256)
				ft.Claims["username"] = "test"
				ft.Claims["admin"] = true

				c.Set("user", ft)
				c.SetPath("/users/")

				Convey("It should error with 401 bad request", func() {
					err := createUserHandler(c).(*echo.HTTPError)
					So(err, ShouldNotBeNil)
					So(err.Code, ShouldEqual, 400)
				})
			})

		})

		Convey("When deleting a user", func() {
			deleteUserSubcriber()

			e := echo.New()
			req := http.Request{Method: "DELETE"}
			rec := httptest.NewRecorder()
			c := e.NewContext(standard.NewRequest(&req, e.Logger()), standard.NewResponse(rec, e.Logger()))

			c.SetPath("/users/:user")
			c.SetParamNames("user")
			c.SetParamValues("1")

			Convey("It should delete the user and return a 200 ok", func() {
				err := deleteUserHandler(c)
				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, 200)
			})
		})

	})
}
