/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"

	"github.com/ernestio/api-gateway/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	log.Println("starting gateway")
	setup()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/auth", controllers.Authenticate)
	e.GET("/status", controllers.GetStatusHandler)

	// Setup JWT auth & protected routes
	api := e.Group("/api")
	api.Use(middleware.JWT([]byte(controllers.Secret)))
	setupRoutes(api)

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
