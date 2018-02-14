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

// Create : Creates an environment build
func Create(au models.User, definition *definition.Definition, raw []byte, dry string) (int, []byte) {
	var e models.Env
	var m models.Mapping

	err := e.FindByName(definition.FullName())
	if err != nil {
		h.L.Error(err.Error())
		return 404, []byte("Environment not found")
	}

	if st, _ := h.IsAuthorizedToResource(&au, h.UpdateEnv, e.GetType(), e.Name); st != 200 {
		// Try submission
		return Submission(au, &e, definition, raw, dry)
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

	err = h.Licensed()
	if err == nil {
		res, err := m.Validate(e.Name)
		if err != nil {
			h.L.Error(err.Error())
			return 400, []byte("build validation failed")
		}

		if res != nil {
			ok := res.Pass()
			if ok != true {
				data, err := json.Marshal(res)
				if err != nil {
					h.L.Error(err.Error())
					return 400, []byte("failed to encode build validation")
				}

				h.L.Error(errors.New("build validation failed"))
				return 400, data
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
		return 500, []byte(`"Couldn't create the build"`)
	}

	if err := b.RequestCreation(&m); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't call build.create"`)
	}

	return http.StatusOK, []byte(`{"id":"` + b.ID + `"}`)
}
