/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Definition : responds to GET /builds/:build/definition with the
// definition of an existing build
func Definition(au models.User, id string) (int, []byte) {
	var err error
	var body []byte
	var b models.Build
	var e models.Env

	if err = b.FindByID(id); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	err = e.FindByID(b.EnvironmentID)
	if err != nil {
		h.L.Error(err.Error())
		return 404, []byte("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	if body, err = b.GetDefinition(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, body
}
