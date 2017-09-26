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
	var list []models.Build
	var body []byte

	err := b.FindByEnvironmentName(env, &list)
	if err != nil {
		h.L.Warning(err.Error())
		return 404, []byte("Environment not found")
	}

	body, err = json.Marshal(list)
	if err != nil {
		return 500, []byte("Internal error")
	}

	return http.StatusOK, body
}
