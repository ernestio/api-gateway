/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package envs

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /services/ with a list of all
// services for current user group
func List(au models.User, project *string) (int, []byte) {
	var body []byte
	var err error
	var r models.Role
	var p models.Project

	query := make(map[string]interface{}, 0)

	if project != nil {
		if !models.IsAlphaNumeric(*project) {
			return 404, models.NewJSONError("Project name contains invalid characters")
		}

		p, err = au.ProjectByName(*project)
		if err != nil {
			h.L.Warning(err.Error())
			return 404, models.NewJSONError("Project not found")
		}

		query["project_id"] = p.ID
	}

	envs, err := au.EnvsBy(query)
	if err != nil {
		h.L.Warning(err.Error())
		return 404, models.NewJSONError("Environment not found")
	}

	pcache := make(map[string][]models.Role)

	for i := range envs {
		var roles []models.Role

		computedRoles := make(map[string]models.Role, 0)

		if pcache[envs[i].Project] == nil {
			var pRoles []models.Role

			err = r.FindAllByResource(envs[i].GetProject(), p.GetType(), &pRoles)
			if err == nil {
				pcache[envs[i].Project] = pRoles
			}
		}

		for _, v := range pcache[envs[i].Project] {
			computedRoles[v.UserID] = v
		}

		err = envs[i].Redact()
		if err != nil {
			return 500, models.NewJSONError(err.Error())
		}

		if err = r.FindAllByResource(envs[i].GetID(), envs[i].GetType(), &roles); err == nil {
			for _, v := range roles {
				computedRoles[v.UserID] = v
			}
		}

		for _, v := range computedRoles {
			envs[i].Members = append(envs[i].Members, v)
		}

		envs[i].Project, envs[i].Name = getProjectEnv(envs[i].Name)
	}

	body, err = json.Marshal(envs)
	if err != nil {
		return 500, models.NewJSONError("Internal error")
	}

	return http.StatusOK, body
}
