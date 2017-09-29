/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/notifications"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetNotificationsHandler : responds to GET /notifications/ with a list of all
// notifications
func GetNotificationsHandler(c echo.Context) (err error) {
	return genericList(c, "notification", notifications.List)
}

// CreateNotificationHandler : responds to POST /notifications/ by creating a notification
// on the data store
func CreateNotificationHandler(c echo.Context) (err error) {
	return genericCreate(c, "notification", notifications.Create)
}

// DeleteNotificationHandler : responds to DELETE /notifications/:id: by deleting an
// existing notification
func DeleteNotificationHandler(c echo.Context) (err error) {
	return genericDelete(c, "notification", notifications.Delete)
}

// UpdateNotificationHandler : ...
func UpdateNotificationHandler(c echo.Context) (err error) {
	return genericUpdate(c, "notification", notifications.Update)
}

// AddEnvToNotificationHandler : ...
func AddEnvToNotificationHandler(c echo.Context) (err error) {
	return entityToNotification(c, notifications.AddEnv, "notifications/add_env")
}

// RmEnvToNotificationHandler : ...
func RmEnvToNotificationHandler(c echo.Context) (err error) {
	return entityToNotification(c, notifications.RmEnv, "notifications/rm_env")
}

// AddProjectToNotificationHandler : ...
func AddProjectToNotificationHandler(c echo.Context) (err error) {
	return entityToNotification(c, notifications.AddEnv, "notifications/add_project")
}

// RmProjectToNotificationHandler : ...
func RmProjectToNotificationHandler(c echo.Context) (err error) {
	return entityToNotification(c, notifications.RmEnv, "notifications/rm_project")
}

type attachEntity func(models.User, string, string) (int, []byte)

func entityToNotification(c echo.Context, fn attachEntity, path string) (err error) {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, path)
	if st != 200 {
		return c.JSONBlob(st, b)
	}

	name := c.Param("notification")
	entity := c.Param("project")
	if env := c.Param("env"); env != "" {
		entity = entity + models.EnvNameSeparator + env
	}

	st, b = fn(au, name, entity)

	return c.JSONBlob(st, b)
}
