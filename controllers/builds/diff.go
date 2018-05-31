/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Diff : Diffs an environment
func Diff(au models.User, env string, request *models.Diff) (int, []byte) {
	var e models.Env
	var m models.Mapping

	if !models.IsAlphaNumeric(env) {
		return 404, models.NewJSONError("Environment name contains invalid characters")
	}

	err := e.FindByName(env)
	if err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	err = m.Diff(env, request.FromID, request.ToID)
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Couldn't diff the specified builds")
	}

	data, err := m.ChangelogJSON()
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Couldn't marshal response")
	}

	return http.StatusOK, data
}
