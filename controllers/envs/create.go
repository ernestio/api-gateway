/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package envs

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /projects/ by creating a
// project on the data store
func Create(au models.User, project string, body []byte) (int, []byte) {
	var err error
	var e models.Env
	var p models.Project
	var existing models.Env

	if e.Map(body) != nil {
		return 400, []byte("Input is not valid")
	}

	err = e.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	err = p.FindByName(project)
	if err != nil {
		return 404, []byte("Specified project does not exist")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, p.GetType(), p.Name); st != 200 {
		return st, res
	}

	e.Name = project + models.EnvNameSeparator + e.Name
	if err := existing.FindByName(e.Name); err == nil {
		return 409, []byte("Specified environment already exists")
	}

	e.ProjectID = p.ID
	e.Type = p.Type

	if err = e.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	if err := au.SetOwner(&e); err != nil {
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(e); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
