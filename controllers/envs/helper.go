/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package envs

import (
	"encoding/json"
	"errors"

	"github.com/ernestio/api-gateway/models"
	"github.com/nu7hatch/gouuid"
)

func getEnvRaw(au models.User, name string) ([]byte, error) {
	filters := make(map[string]interface{}, 0)
	filters["name"] = name

	ss, err := au.EnvsBy(filters)
	if err != nil {
		return nil, err
	}

	if len(ss) == 0 {
		return nil, errors.New("Not found")
	}

	body, err := json.Marshal(ss[0])
	if err != nil {
		return nil, errors.New("Internal error")
	}
	return body, nil
}

// Generates an environment id composed by a random uuid, and
// a valid generated stream id
func generateEnvID(salt string) string {
	id, _ := uuid.NewV4()

	return id.String()
}
