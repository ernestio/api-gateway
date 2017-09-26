/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Import : Imports an environment
func Import(au models.User, env string, filters []string) (int, []byte) {
	var e models.Env
	var m models.Mapping

	err := e.FindByName(env)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't get the environment"`)
	}

	err = m.Import(env, filters)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the import build"`)
	}

	b := models.Build{
		ID:            m["id"].(string),
		EnvironmentID: e.ID,
		UserID:        au.ID,
		Username:      au.Username,
		Type:          "import",
		Mapping:       m,
	}

	err = b.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't create the build"`)
	}

	if err := b.RequestImport(&m); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't call build.import"`)
	}

	return http.StatusOK, []byte(`{"id":"` + b.ID + `"}`)
}
