/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"os"
	"time"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/labstack/echo"
)

func setup() {
	n = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()

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
	u.POST("/", createUserHandler)
	u.PUT("/:user", updateUserHandler)
	u.DELETE("/:user", deleteUserHandler)

	// Setup group routes
	g := api.Group("/groups")
	g.GET("/", getGroupsHandler)
	g.GET("/:group", getGroupHandler)
	g.POST("/", createGroupHandler)
	g.PUT("/:group", updateGroupHandler)
	g.DELETE("/:group", deleteGroupHandler)
	g.POST("/:group/users/", addUserToGroupHandler)
	g.DELETE("/:group/users/:user", deleteUserFromGroupHandler)
	g.POST("/:group/datacenters/", addDatacenterToGroupHandler)
	g.DELETE("/:group/datacenters/:datacenter", deleteDatacenterFromGroupHandler)

	// Setup datacenter routes
	d := api.Group("/datacenters")
	d.GET("/", getDatacentersHandler)
	d.GET("/:datacenter", getDatacenterHandler)
	d.POST("/", createDatacenterHandler)
	d.PUT("/:datacenter", updateDatacenterHandler)
	d.DELETE("/:datacenter", deleteDatacenterHandler)

	// Setup logger routes
	l := api.Group("/loggers")
	l.GET("/", getLoggersHandler)
	l.POST("/", createLoggerHandler)
	l.DELETE("/:logger", deleteLoggerHandler)

	// Setup service routes
	s := api.Group("/services")
	s.GET("/", getServicesHandler)
	s.GET("/:service", getServiceHandler)
	s.GET("/search/", searchServicesHandler)
	s.GET("/:service/builds/", getServiceBuildsHandler)
	s.GET("/:service/builds/:build", getServiceBuildHandler)
	s.POST("/", createServiceHandler)
	s.POST("/import/", createServiceHandler)
	s.POST("/uuid/", createUUIDHandler)
	s.POST("/:service/reset/", resetServiceHandler)
	s.PUT("/:service", updateServiceHandler)
	s.DELETE("/:name", deleteServiceHandler)
	s.DELETE("/:name/force/", forceServiceDeletionHandler)

	// Setup components
	comp := api.Group("/components")
	comp.GET("/nats/", getAllComponentsHandler)
	comp.GET("/network/", getAllComponentsHandler)
	comp.GET("/route53/", getAllComponentsHandler)
	comp.GET("/s3/", getAllComponentsHandler)
	comp.GET("/elb/", getAllComponentsHandler)
	comp.GET("/vpc/", getAllComponentsHandler)
	comp.GET("/instance/", getAllComponentsHandler)
	comp.GET("/firewall/", getAllComponentsHandler)
	comp.GET("/ebs_volume/", getAllComponentsHandler)
	comp.GET("/rds_cluster/", getAllComponentsHandler)
	comp.GET("/rds_instance/", getAllComponentsHandler)

	// Setup reports
	rep := api.Group("/reports")
	rep.GET("/usage/", getUsageReportHandler)

}
