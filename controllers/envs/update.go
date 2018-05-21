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

// Update : responds to PUT /projects/:project:/envs/:env/ by updating an
// existing environment
func Update(au models.User, name string, body []byte) (int, []byte) {
	var err error
	var resp []byte
	var e models.Env
	var input models.Env

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateEnv, input.GetType(), name); st != 200 {
		return st, res
	}

	if err = json.Unmarshal(body, &input); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Get existing environment
	if err := e.FindByName(name); err != nil {
		return 404, models.NewJSONError(err.Error())
	}

	e.Options = input.Options
	e.Schedules = input.Schedules
	e.Credentials = input.Credentials

	if err := e.Save(); err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	for _, r := range input.Members {
		for _, er := range e.Members {
			// create role
			if r.ID == 0 {
				err = r.Save()
				if err != nil {
					h.L.Error(err.Error())
					return http.StatusBadRequest, models.NewJSONError(err.Error())
				}
			}

			// update role
			if r.ID == er.ID && r.Role != er.Role {
				err = r.Save()
				if err != nil {
					h.L.Error(err.Error())
					return http.StatusBadRequest, models.NewJSONError(err.Error())
				}
			}
		}
	}

	for _, er := range e.Members {
		var exists bool

		for _, r := range input.Members {
			if r.ID == er.ID {
				exists = true
			}
		}

		// delete roles
		if !exists {
			err = er.Delete()
			if err != nil {
				h.L.Error(err.Error())
				return http.StatusBadRequest, models.NewJSONError(err.Error())
			}
		}
	}

	e.Members = input.Members

	resp, err = json.Marshal(e)
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, models.NewJSONError(err.Error())
	}

	return http.StatusOK, resp
}
