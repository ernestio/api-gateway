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

	query := make(map[string]interface{}, 0)

	if project != nil {
		p, err := au.ProjectByName(*project)
		if err != nil {
			h.L.Warning(err.Error())
			return 404, []byte("Project not found")
		}

		query["project_id"] = p.ID
	}

	envs, err := au.EnvsBy(query)
	if err != nil {
		h.L.Warning(err.Error())
		return 404, []byte("Environment not found")
	}

	for i := range envs {
		envs[i].Project, envs[i].Name = getProjectEnv(envs[i].Name)
	}

	body, err = json.Marshal(envs)
	if err != nil {
		return 500, []byte("Internal error")
	}

	return http.StatusOK, body
}
