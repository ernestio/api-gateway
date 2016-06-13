package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/nats-io/nats"
)

var n *nats.Conn
var secret string

func main() {
	log.Println("starting gateway")
	setup()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Post("/auth", authenticate)

	// Setup JWT auth & protected routes
	api := e.Group("/api")
	api.Use(middleware.JWT([]byte(secret)))
	setupRoutes(api)

	e.Run(standard.New(":8080"))
}
