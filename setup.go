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
	// Setup user routes
	api.GET("/users/", getUsersHandler)
	api.GET("/users/:user", getUserHandler)
	api.Post("/users/", createUserHandler)
	api.Put("/users/:user", updateUserHandler)
	api.Delete("/users/:user", deleteUserHandler)

	// Setup group routes
	api.GET("/groups/", getGroupsHandler)
	api.GET("/groups/:group", getGroupHandler)
	api.Post("/groups/", createGroupHandler)
	api.Put("/groups/:group", updateGroupHandler)
	api.Delete("/groups/:group", deleteGroupHandler)

	// Setup datacenter routes
	api.GET("/datacenters/", getDatacentersHandler)
	api.GET("/datacenters/:datacenter", getDatacenterHandler)
	api.Post("/datacenters/", createDatacenterHandler)
	api.Put("/datacenters/:datacenter", updateDatacenterHandler)
	api.Delete("/datacenters/:datacenter", deleteDatacenterHandler)

	// Setup service routes
	api.GET("/services/", getServicesHandler)
	api.GET("/services/:service", getServiceHandler)
	api.GET("/services/search/", searchServicesHandler)
	api.GET("/services/:service/builds/", getServiceBuildsHandler)
	api.GET("/services/:service/builds/:build", getServiceBuildHandler)
	api.Post("/services/", createServiceHandler)
	api.Post("/services/uuid/", createUUIDHandler)
	api.Post("/services/:service/reset/", resetServiceHandler)
	api.Put("/services/:service", updateServiceHandler)
	api.Delete("/services/:service", deleteServiceHandler)
}
