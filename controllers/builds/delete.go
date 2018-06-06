/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : Deletes an environment by name, generating a delete build
func Delete(au models.User, name string) (int, []byte) {
	var e models.Env
	var m models.Mapping

	if !models.IsAlphaNumeric(name) {
		return 404, models.NewJSONError("Environment name contains invalid characters")
	}

	err := e.FindByName(name)
	if err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.DeleteEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	err = m.Delete(name, au)
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Couldn't map the environment")
	}

	b := models.Build{
		ID:            m["id"].(string),
		EnvironmentID: e.ID,
		UserID:        au.ID,
		Username:      au.Username,
		Type:          "destroy",
		Mapping:       m,
	}

	err = b.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 400, models.NewJSONError("Environment is already applying some changes, please wait until they are done")
	}

	if err := b.RequestDeletion(&m); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Couldn't call build.delete")
	}

	return http.StatusOK, []byte(`{"id":"` + b.ID + `"}`)
}
