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

// Get : responds to GET /policies/:name/revisions/revision with the policy
// details
func GetDocument(au models.User, name, revision string) (int, []byte) {
	var err error
	var body []byte
	var policy models.Policy
	var document models.PolicyDocument

	if err = policy.GetByName(name, &policy); err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("policy not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetPolicy, policy.GetType(), policy.GetID()); st != 200 {
		return st, res
	}

	if err = document.GetByRevision(name, revision, &document); err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("policy revision not found")
	}

	if body, err = json.Marshal(document); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal server error")
	}
	return http.StatusOK, body
}
