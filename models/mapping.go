/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ernestio/mapping"
	"github.com/ernestio/mapping/definition"
)

type policyResponse struct {
	status string `json:"status"`
	output string `json:"output"`
}

// Mapping : graph mapping
type Mapping map[string]interface{}

// Apply : apply a definition
func (m *Mapping) Apply(d *definition.Definition, au User) error {
	return m.apply(d, au, false)
}

// Submission : submit a definition
func (m *Mapping) Submission(d *definition.Definition, au User) error {
	return m.apply(d, au, true)
}

// Apply : apply a definition
func (m *Mapping) apply(d *definition.Definition, au User, changelog bool) error {
	mr := mapping.New(N, d.FullName())

	mr.Changelog = changelog

	err := mr.Apply(d)
	if err != nil {
		return err
	}

	mr.Result["user_id"] = au.ID
	mr.Result["username"] = au.Username

	*m = mr.Result

	return nil
}

// Delete : get mapping for deleting an environment
func (m *Mapping) Delete(env string, au User) error {
	mr := mapping.New(N, env)

	err := mr.Delete()
	if err != nil {
		return err
	}

	mr.Result["user_id"] = au.ID
	mr.Result["username"] = au.Username

	*m = mr.Result

	return nil
}

// Import : get mapping for importing an environment
func (m *Mapping) Import(env string, filters []string, au User) error {
	mr := mapping.New(N, env)

	err := mr.Import(filters)
	if err != nil {
		return err
	}

	mr.Result["user_id"] = au.ID
	mr.Result["username"] = au.Username

	*m = mr.Result

	return nil
}

// Diff : diff two builds by id
func (m *Mapping) Diff(env, from, to string) error {
	mr := mapping.New(N, env)

	err := mr.Diff(from, to)
	if err != nil {
		return err
	}

	*m = mr.Result

	return nil
}

// Validate : checks a map against any attached policies.
func (m *Mapping) Validate(env string) (string, error) {
	req := fmt.Sprintf(`{"environment": "%s", "mapping": "%s"}`, env, m)
	msg, err := N.Request("policy.check", []byte(req), 1*time.Second)
	if err != nil {
		return "", err
	}

	var pr policyResponse
	err = json.Unmarshal(msg.Data, &pr)
	if err != nil {
		return "", err
	}

	if pr.status == "error" {
		return pr.output, errors.New("policy error")
	}

	return pr.output, nil
}

// Changelog : returns the mappings changelog if present
func (m *Mapping) ChangelogJSON() ([]byte, error) {
	if (*m)["changelog"] == nil {
		return json.Marshal([]string{})
	}

	return json.Marshal((*m)["changelog"])
}

// ToJSON : serializes the mapping to json
func (m *Mapping) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
