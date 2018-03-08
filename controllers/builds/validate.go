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

// Validate : Validate an environment
func Validate(au models.User, env string, action *models.Action) (int, []byte) {
	var e models.Env
	var b models.Build
	var m models.Mapping
	var builds []models.Build

	err := e.FindByName(env)
	if err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	err = b.FindByEnvironmentName(env, &builds)
	if err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	if len(builds) < 1 {
		return 400, models.NewJSONError("could not validate environment, there are no builds.")
	}

	m, err = builds[0].GetRawMapping()
	if err != nil {
		return 400, models.NewJSONError("unable to get environment build")
	}

	validation, err := m.Validate(env)
	if err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	data, err := json.Marshal(validation)
	if err != nil {
		return 500, models.NewJSONError("could not process sync request")
	}

	return http.StatusOK, data
}
