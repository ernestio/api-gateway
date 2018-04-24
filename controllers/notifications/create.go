/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /notifications/ by creating a notification
// on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var n models.Notification
	var existing models.Notification
	var e models.Env
	var p models.Project
	var envs []models.Env
	var projects []models.Project

	if n.Map(body) != nil {
		return http.StatusBadRequest, models.NewJSONError("Invalid input")
	}

	err := existing.FindByName(n.Name, &existing)
	if err == nil {
		return 409, models.NewJSONError("Specified notifiation already exists")
	}

	err = n.Validate()
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

	err = p.FindAll(au, &projects)
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

	err = e.FindAll(au, &envs)
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

SOURCELOOP:
	for _, source := range n.Sources {
		for _, project := range projects {
			if project.Name == source {
				continue SOURCELOOP
			}
		}

		for _, env := range envs {
			if env.Name == source {
				continue SOURCELOOP
			}
		}

		return 400, models.NewJSONError("notification source '" + source + "' does not exist")
	}

	err = n.Save()
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

	body, err = json.Marshal(n)
	if err != nil {
		return 500, models.NewJSONError("Internal error")
	}

	return http.StatusOK, body
}
