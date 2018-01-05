/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package policies

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Get : responds to GET /policies/:name/ with the policy
// details
func Get(au models.User, name string) (int, []byte) {
	var err error
	var body []byte
	var policy models.Policy

	if err = policy.FindByName(name, &policy); err != nil {
		h.L.Error(err.Error())
		return 404, []byte("Policy not found")
	}

	if body, err = json.Marshal(policy); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal serveier error")
	}
	return http.StatusOK, body
}
