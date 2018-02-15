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

// List : responds to GET /builds/ with a list of all
// builds for current user group
func List(au models.User, env string) (int, []byte) {
	var b models.Build
	var e models.Env
	var list []models.Build
	var body []byte

	err := e.FindByName(env)
	if err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	err = b.FindByEnvironmentName(env, &list)
	if err != nil {
		h.L.Warning(err.Error())
		return 404, models.NewJSONError("Build not found")
	}

	for i := len(list) - 1; i >= 0; i-- {
		if list[i].Type == "sync" && list[i].Status == "done" {
			list = append(list[:i], list[i+1:]...)
		}
	}

	body, err = json.Marshal(list)
	if err != nil {
		return 500, models.NewJSONError("Internal error")
	}

	return http.StatusOK, body
}
