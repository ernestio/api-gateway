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
func Submission(au models.User, env *models.Env, definition *definition.Definition, raw []byte, dry string) (int, []byte) {
	var e models.Env
	var m models.Mapping

	err := e.FindByName(definition.FullName())
	if err != nil {
		h.L.Error(err.Error())
		return 404, []byte("Environment not found")
	}

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

	err = m.Apply(definition, au)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the environment"`)
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
