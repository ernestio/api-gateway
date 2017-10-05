/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"

	"github.com/ernestio/mapping"
	"github.com/ernestio/mapping/definition"
)

// Mapping : graph mapping
type Mapping map[string]interface{}

// Apply : apply a definition
func (m *Mapping) Apply(d *definition.Definition) error {
	var err error

	mr := mapping.New(N, d.FullName())
	err = mr.Apply(d)
	if err != nil {
		return err
	}

	*m = mr.Result

	return nil
}

// Delete : get mapping for deleting an environment
func (m *Mapping) Delete(env string) error {
	var err error

	mr := mapping.New(N, env)
	err = mr.Delete()
	if err != nil {
		return err
	}

	*m = mr.Result

	return nil
}

// Import : get mapping for importing an environment
func (m *Mapping) Import(env string, filters []string) error {
	var err error

	mr := mapping.New(N, env)
	err = mr.Import(filters)
	if err != nil {
		return err
	}

	*m = mr.Result

	return nil
}

// ToJSON : serializes the mapping to json
func (m *Mapping) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
