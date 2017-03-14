/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"github.com/ernestio/api-gateway/controllers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Route : set up router and starts the server
func Route() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	L.Info("Setting up root routes")
	setupRoot(e)
	L.Info("Setting up api routes")
	setupAPI(e)
	L.Info("Starting server")
	start(e)

	return e
}

func setupRoot(e *echo.Echo) {
	e.POST("/auth", controllers.Authenticate)
	e.GET("/status", controllers.GetStatusHandler)
}

func setupAPI(e *echo.Echo) {
	api := e.Group("/api")
	api.Use(middleware.JWT([]byte(controllers.Secret)))

	ss := api.Group("/session")
	ss.GET("/", controllers.GetSessionsHandler)

	// Setup user routes
	u := api.Group("/users")
	u.GET("/", controllers.GetUsersHandler)
	u.GET("/:user", controllers.GetUserHandler)
	u.POST("/", controllers.CreateUserHandler)
	u.PUT("/:user", controllers.UpdateUserHandler)
	u.DELETE("/:user", controllers.DeleteUserHandler)

	// Setup group routes
	g := api.Group("/groups")
	g.GET("/", controllers.GetGroupsHandler)
	g.GET("/:group", controllers.GetGroupHandler)
	g.POST("/", controllers.CreateGroupHandler)
	g.PUT("/:group", controllers.UpdateGroupHandler)
	g.DELETE("/:group", controllers.DeleteGroupHandler)
	g.POST("/:group/users/", controllers.AddUserToGroupHandler)
	g.DELETE("/:group/users/:user", controllers.DeleteUserFromGroupHandler)
	g.POST("/:group/datacenters/", controllers.AddDatacenterToGroupHandler)
	g.DELETE("/:group/datacenters/:datacenter", controllers.DeleteDatacenterFromGroupHandler)

	// Setup datacenter routes
	d := api.Group("/datacenters")
	d.GET("/", controllers.GetDatacentersHandler)
	d.GET("/:datacenter", controllers.GetDatacenterHandler)
	d.POST("/", controllers.CreateDatacenterHandler)
	d.PUT("/:datacenter", controllers.UpdateDatacenterHandler)
	d.DELETE("/:datacenter", controllers.DeleteDatacenterHandler)

	// Setup logger routes
	l := api.Group("/loggers")
	l.GET("/", controllers.GetLoggersHandler)
	l.POST("/", controllers.CreateLoggerHandler)
	l.DELETE("/:logger", controllers.DeleteLoggerHandler)

	// Setup service routes
	s := api.Group("/services")
	s.GET("/", controllers.GetServicesHandler)
	s.GET("/:service", controllers.GetServiceHandler)
	s.GET("/search/", controllers.SearchServicesHandler)
	s.GET("/:service/builds/", controllers.GetServiceBuildsHandler)
	s.GET("/:service/builds/:build", controllers.GetServiceBuildHandler)
	s.POST("/", controllers.CreateServiceHandler)
	s.POST("/import/", controllers.CreateServiceHandler)
	s.POST("/uuid/", controllers.CreateUUIDHandler)
	s.POST("/:service/reset/", controllers.ResetServiceHandler)
	s.PUT("/:service", controllers.UpdateServiceHandler)
	s.DELETE("/:name", controllers.DeleteServiceHandler)
	s.DELETE("/:name/force/", controllers.ForceServiceDeletionHandler)

	// Setup components
	comp := api.Group("/components")
	comp.GET("/nats/", controllers.GetAllComponentsHandler)
	comp.GET("/network/", controllers.GetAllComponentsHandler)
	comp.GET("/route53/", controllers.GetAllComponentsHandler)
	comp.GET("/s3/", controllers.GetAllComponentsHandler)
	comp.GET("/elb/", controllers.GetAllComponentsHandler)
	comp.GET("/vpc/", controllers.GetAllComponentsHandler)
	comp.GET("/instance/", controllers.GetAllComponentsHandler)
	comp.GET("/firewall/", controllers.GetAllComponentsHandler)
	comp.GET("/ebs_volume/", controllers.GetAllComponentsHandler)
	comp.GET("/rds_cluster/", controllers.GetAllComponentsHandler)
	comp.GET("/rds_instance/", controllers.GetAllComponentsHandler)

	// Setup reports
	rep := api.Group("/reports")
	rep.GET("/usage/", controllers.GetUsageReportHandler)

}

func start(e *echo.Echo) {
	c := models.Config{}
	port := c.GetServerPort()
	e.Logger.Fatal(e.Start(":" + port))
}
