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

// Sync : Syncs an environment
func Sync(au models.User, env string, action *models.Action) (int, []byte) {
	var e models.Env

	if !models.IsAlphaNumeric(env) {
		return 404, models.NewJSONError("Environment name contains invalid characters")
	}

	err := e.FindByName(env)
	if err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	id, err := e.RequestSync(au)
	if err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	action.ResourceType = "build"
	action.ResourceID = id
	action.Status = "syncing"

	data, err := json.Marshal(action)
	if err != nil {
		return 500, models.NewJSONError("could not process sync request")
	}

	return http.StatusOK, data
}
