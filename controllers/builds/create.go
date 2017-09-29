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

// Create : Creates an environment build
func Create(au models.User, definition *definition.Definition, raw []byte, dry string) (int, []byte) {
	var e models.Env
	var m models.Mapping

	err := e.FindByName(definition.FullName())
	if err != nil {
		h.L.Error(err.Error())
		return 404, []byte("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	err = m.Apply(definition)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the environment"`)
	}

	if dry == "true" {
		res, err := views.RenderDefinition(m)
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
