package main

import (
	"log"
	"os"

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
		panic("No JWT secret was set!")
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
}
