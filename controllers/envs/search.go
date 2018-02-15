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

// Search : Finds all services
func Search(au models.User, query map[string]interface{}) (int, []byte) {
	envs, err := au.EnvsBy(query)
	if err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	b, err := json.Marshal(envs)
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal error")
	}

	return http.StatusOK, b
}
