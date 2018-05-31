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
	"github.com/ernestio/mapping/validation"
)

// Create : Creates an environment build
func Create(au models.User, definition *definition.Definition, raw []byte, dry string) (int, []byte) {
	var e models.Env
	var m models.Mapping
	var validation *validation.Validation

	if !models.IsAlphaNumeric(definition.FullName()) {
		return 404, models.NewJSONError("Notification name contains invalid characters")
	}

	err := e.FindByName(definition.FullName())
	if err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Environment not found")
	}

	if st, _ := h.IsAuthorizedToResource(&au, h.UpdateEnv, e.GetType(), e.Name); st != 200 {
		// Try submission
		return Submission(au, &e, definition, raw, dry)
	}

	err = m.Apply(definition, au)
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError(err.Error())
	}

	if dry == "true" {
		res, err := views.RenderChanges(m)
		if err != nil {
			h.L.Error(err.Error())
			return 400, models.NewJSONError("Internal error")
		}
		return http.StatusOK, res
	}

	if h.Licensed() == nil {
		validation, err = m.Validate(e.Name)
		if err != nil {
			h.L.Error(err.Error())
			return 400, models.NewJSONError("could not validate build")
		}

		if validation != nil {
			if validation.Passed() != true {
				h.L.Error(errors.New("build validation failed"))
				return 400, models.NewJSONValidationError("build validation failed", validation)
			}
		}
	}

	b := models.Build{
		ID:            m["id"].(string),
		EnvironmentID: e.ID,
		UserID:        au.ID,
		Username:      au.Username,
		Type:          "apply",
		Mapping:       m,
		Definition:    string(raw),
	}

	err = b.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Couldn't create the build")
	}

	if err := b.RequestCreation(&m); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Couldn't call build.create")
	}

	br := models.BuildDetails{
		ID:         b.ID,
		Status:     b.Status,
		Validation: validation,
	}

	data, err := json.Marshal(br)
	if err != nil {
		h.L.Error(err.Error())
		return 400, models.NewJSONError("Internal error")
	}

	return http.StatusOK, data
}
