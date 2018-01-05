/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
	"github.com/ernestio/mapping/definition"
)

// Submission : Submits an environment build for approval
func Submission(au models.User, e *models.Env, definition *definition.Definition, raw []byte, dry string) (int, []byte) {
	var m models.Mapping

	submissions, _ := e.Options["submissions"].(bool)
	if !submissions {
		return 403, h.AuthNonOwner
	}

	if st, res := h.IsLicensed(&au, h.SubmitBuild); st != 200 {
		return st, res
	}

	if st, res := h.IsAuthorizedToReadResource(&au, h.SubmitBuild, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	err := m.Apply(definition, au)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the environment"`)
	}

	changes, ok := m["changes"].([]interface{})
	if !ok || len(changes) < 1 {
		h.L.Error(err.Error())
		return 400, []byte(`"The provided definition contains no changes."`)
	}

	if dry == "true" {
		res, err := views.RenderChanges(m)
		if err != nil {
			h.L.Error(err.Error())
			return 400, []byte("Internal error")
		}
		return http.StatusOK, res
	}

	b := models.Build{
		ID:            m["id"].(string),
		EnvironmentID: e.ID,
		UserID:        au.ID,
		Username:      au.Username,
		Type:          "submission",
		Mapping:       m,
		Definition:    string(raw),
	}

	err = b.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't create the build"`)
	}

	return http.StatusOK, []byte(`{"id":"` + b.ID + `", "status":"submitted"}`)
}
