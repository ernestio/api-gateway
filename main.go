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
}

func main() {
	log.Println("starting")
	setup()

	e := echo.New()
	e.Use(middleware.JWT([]byte("test")))

	// Setup users routes
	e.GET("/users/", getUsersHander)
	e.GET("/users/:user", getUserHander)
	e.Post("/users/", createUserHander)
	e.Put("/users/:user", updateUserHander)
	e.Delete("/users/:user", deleteUserHander)

	e.Run(standard.New(":8080"))
}
