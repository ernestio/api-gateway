/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package roles

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Get : responds to GET /roles/:name/ with the role
// details
func Get(au models.User, id string) (int, []byte) {
	var err error
	var body []byte
	var role models.Role

	if !models.IsAlphaNumeric(id) {
		return 404, models.NewJSONError("Role ID contains invalid characters")
	}

	if err = role.FindByID(id, &role); err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Role not found")
	}

	if body, err = json.Marshal(role); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal serveier error")
	}
	return http.StatusOK, body
}
