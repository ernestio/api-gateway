/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/nats-io/nats"
)

func setup() {
	var err error
	natsURI := os.Getenv("NATS_URI")
	if natsURI == "" {
		natsURI = nats.DefaultURL
	}

	n, err = nats.Connect(natsURI)
	if err != nil {
		log.Panic(err)
	}

	secret = os.Getenv("JWT_SECRET")
	if secret == "" {
		token, err := n.Request("config.get.jwt_token", []byte(""), 1*time.Second)
		if err != nil {
			panic("Can't get jwt_config config")
		}

		secret = string(token.Data)
	}
}

func setupRoutes(api *echo.Group) {
	// Setup session routes
	ss := api.Group("/session")
	ss.GET("/", getSessionsHandler)

	// Setup user routes
	u := api.Group("/users")
	u.GET("/", getUsersHandler)
	u.GET("/:user", getUserHandler)
	u.Post("/", createUserHandler)
	u.Put("/:user", updateUserHandler)
	u.Delete("/:user", deleteUserHandler)

	// Setup group routes
	g := api.Group("/groups")
	g.GET("/", getGroupsHandler)
	g.GET("/:group", getGroupHandler)
	g.Post("/", createGroupHandler)
	g.Put("/:group", updateGroupHandler)
	g.Delete("/:group", deleteGroupHandler)

	// Setup datacenter routes
	d := api.Group("/datacenters")
	d.GET("/", getDatacentersHandler)
	d.GET("/:datacenter", getDatacenterHandler)
	d.Post("/", createDatacenterHandler)
	d.Put("/:datacenter", updateDatacenterHandler)
	d.Delete("/:datacenter", deleteDatacenterHandler)

	// Setup service routes
	s := api.Group("/services")
	s.GET("/", getServicesHandler)
	s.GET("/:service", getServiceHandler)
	s.GET("/search/", searchServicesHandler)
	s.GET("/:service/builds/", getServiceBuildsHandler)
	s.GET("/:service/builds/:build", getServiceBuildHandler)
	s.Post("/", createServiceHandler)
	s.Post("/uuid/", createUUIDHandler)
	s.Post("/:service/reset/", resetServiceHandler)
	s.Put("/:service", updateServiceHandler)
	s.Delete("/:service", deleteServiceHandler)
}
