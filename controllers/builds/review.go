/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Review : Resolves a build that is queued pending approval
func Review(au models.User, env string, action *models.Action) (int, []byte) {
	var e models.Env

	err := e.FindByName(env)
	if err != nil {
		h.L.Error(err.Error())
		return 404, []byte("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	id, err := e.RequestReview(au, action.Options.Resolution)
	if err != nil {
		return 500, []byte(err.Error())
	}

	action.Status = "done"

	if id != "" {
		action.ResourceType = "build"
		action.ResourceID = id
		action.Status = "in_progress"
	}

	data, err := json.Marshal(action)
	if err != nil {
		return 500, []byte("could not process sync request")
	}

	return http.StatusOK, data
}
