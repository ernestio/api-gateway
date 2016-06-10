package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/nats-io/nats"
)

var n *nats.Conn
var secret string

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

func main() {
	log.Println("starting gateway")
	setup()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// login
	e.Post("/auth", authenticate)

	// Setup JWT auth
	api := e.Group("/api")
	api.Use(middleware.JWT([]byte(secret)))

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

	e.Run(standard.New(":8080"))
}
