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

// CreateDocument : responds to POST /policies/:policy/revisions/ by creating a policy
// on the data store
func CreateDocument(au models.User, name string, body []byte) (int, []byte) {
	var policy models.Policy
	var document models.PolicyDocument
	var err error

	if document.Map(body) != nil {
		return http.StatusBadRequest, models.NewJSONError("Invalid input")
	}

	err = document.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError(err.Error())
	}

	if err = policy.GetByName(name, &policy); err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("policy not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.UpdatePolicy, policy.GetType(), policy.GetID()); st != 200 {
		return st, res
	}

	document.Username = au.Username
	document.PolicyID = policy.ID

	err = document.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 400, models.NewJSONError(err.Error())
	}

	if body, err = json.Marshal(document); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}
	return http.StatusOK, body
}
