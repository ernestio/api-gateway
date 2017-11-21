/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"github.com/ernestio/api-gateway/controllers"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Route : set up router and starts the server
func Route() *echo.Echo {
	e := echo.New()
	e.Pre(middleware.AddTrailingSlash())
	//e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h.L.Info("Setting up root routes")
	setupRoot(e)
	h.L.Info("Setting up api routes")
	setupAPI(e)
	h.L.Info("Starting server")
	start(e)

	return e
}

func setupRoot(e *echo.Echo) {
	e.POST("/auth/", controllers.AuthenticateHandler)
	e.GET("/status/", controllers.GetStatusHandler)
}

func setupAPI(e *echo.Echo) {
	api := e.Group("/api")
	api.Use(middleware.JWT([]byte(controllers.Secret)))

	ss := api.Group("/session")
	ss.GET("/", controllers.GetSessionsHandler)

	// Setup user routes
	u := api.Group("/users")
	u.GET("/", controllers.GetUsersHandler)
	u.GET("/:user/", controllers.GetUserHandler)
	u.POST("/", controllers.CreateUserHandler)
	u.PUT("/:user/", controllers.UpdateUserHandler)
	u.DELETE("/:user/", controllers.DeleteUserHandler)

	// Setup roles routes
	r := api.Group("/roles")
	r.GET("/", controllers.GetRolesHandler)
	r.GET("/:role/", controllers.GetRoleHandler)
	r.POST("/", controllers.CreateRoleHandler)
	r.DELETE("/", controllers.DeleteRoleHandler)
	r.DELETE("/:role/", controllers.DeleteRoleByIDHandler)

	// Setup logger routes
	l := api.Group("/loggers")
	l.GET("/", controllers.GetLoggersHandler)
	l.POST("/", controllers.CreateLoggerHandler)
	l.DELETE("/:logger/", controllers.DeleteLoggerHandler)

	// Setup project routes
	d := api.Group("/projects")
	d.GET("/", controllers.GetDatacentersHandler)
	d.GET("/:project/", controllers.GetDatacenterHandler)
	d.POST("/", controllers.CreateDatacenterHandler)
	d.PUT("/:project/", controllers.UpdateDatacenterHandler)
	d.DELETE("/:project/", controllers.DeleteDatacenterHandler)

	// Setup env routes
	d.GET("/:project/envs/", controllers.GetEnvsHandler)
	d.POST("/:project/envs/", controllers.CreateEnvHandler)
	d.PUT("/:project/envs/:env/", controllers.UpdateEnvHandler)
	d.GET("/:project/envs/:env/", controllers.GetEnvHandler)
	d.DELETE("/:project/envs/:env/", controllers.DeleteEnvHandler)

	// Setup build routes
	d.GET("/:project/envs/:env/builds/", controllers.GetBuildsHandler)
	d.POST("/:project/envs/:env/builds/", controllers.CreateBuildHandler)
	d.GET("/:project/envs/:env/builds/:build/", controllers.GetBuildHandler)
	d.GET("/:project/envs/:env/builds/:build/mapping/", controllers.GetBuildMappingHandler)
	d.GET("/:project/envs/:env/builds/:build/definition/", controllers.GetBuildDefinitionHandler)
	d.POST("/:project/envs/:env/actions/", controllers.ActionHandler)

	d.DELETE("/:project/envs/:env/actions/force/", controllers.ForceEnvDeletionHandler)

	// Setup service routes
	s := api.Group("/envs")
	s.GET("/", controllers.GetEnvsHandler)
	s.GET("/search/", controllers.SearchEnvsHandler)

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

	// Setup notifications
	not := api.Group("/notifications")
	not.GET("/", controllers.GetNotificationsHandler)
	not.GET("/:notification/", controllers.GetNotificationHandler)
	not.POST("/", controllers.CreateNotificationHandler)
	not.PUT("/:notification/", controllers.UpdateNotificationHandler)
	not.DELETE("/:notification/", controllers.DeleteNotificationHandler)
	not.POST("/:notification/projects/:project/", controllers.AddProjectToNotificationHandler)
	not.DELETE("/:notification/projects/:project/", controllers.RmProjectToNotificationHandler)
	not.POST("/:notification/projects/:project/envs/:env/", controllers.AddEnvToNotificationHandler)
	not.DELETE("/:notification/projects/:project/envs/:env/", controllers.RmEnvToNotificationHandler)
}

func start(e *echo.Echo) {
	c := models.Config{}
	port := c.GetServerPort()
	e.Logger.Fatal(e.Start(":" + port))
}
