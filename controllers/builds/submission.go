/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"encoding/json"
	"errors"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
	"github.com/ernestio/mapping/definition"
)

// Submission : Submits an environment build for approval
func Submission(au models.User, e *models.Env, definition *definition.Definition, raw []byte, dry string) (int, []byte) {
	var m models.Mapping
	var validation *models.BuildValidateResponse

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

	err := m.Submission(definition, au)
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("could not map the environment")
	}

	if h.Licensed() == nil {
		validation, err = m.Validate(e.Name)
		if err != nil {
			h.L.Error(err.Error())
			return 400, models.NewJSONError("could not validate build")
		}

		if validation.Passed() != true {
			h.L.Error(errors.New("build validation failed"))
			return 400, models.NewJSONValidationError("build validation failed", validation)
		}
	}

	changes, ok := m["changes"].([]interface{})
	if !ok || len(changes) < 1 {
		return 400, models.NewJSONError("The provided definition contains no changes.")
	}

	if dry == "true" {
		res, err := views.RenderChanges(m)
		if err != nil {
			h.L.Error(err.Error())
			return 400, models.NewJSONError("Internal error")
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
		return 500, models.NewJSONError("could not create the build")
	}

	br := models.BuildDetails{
		ID:         b.ID,
		Status:     "submitted",
		Validation: validation,
	}

	data, err := json.Marshal(br)
	if err != nil {
		h.L.Error(err.Error())
		return 400, models.NewJSONError("Internal error")
	}

	return http.StatusOK, data
}
