/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ernestio/mapping"
	"github.com/ernestio/mapping/definition"
)

type BuildValidate struct {
	Mapping  *Mapping `json:"mapping"`
	Policies []Policy `json:"policies"`
}

type BuildValidateResponse struct {
	Version    string     `json:"version"`
	Controls   []Control  `json:"controls"`
	Statistics Statistics `json:"statistics"`
}

type Control struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	Status    string `json:"status"`
	CodeDesc  string `json:"code_desc"`
	Message   string `json:"message"`
}

type Statistics struct {
	Duration float64 `json:"duration"`
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
func (m *Mapping) Validate(project, environment string) (*BuildValidateResponse, error) {
	policyReq := fmt.Sprintf(`{"environment": ["%s/%s"]}`, project, environment)
	msg, err := N.Request("policy.find", []byte(policyReq), 2*time.Second)
	if err != nil {
		return nil, err
	}

	var p []Policy
	err = json.Unmarshal(msg.Data, &p)
	if err != nil {
		return nil, err
	}

	if len(p) == 0 {
		return nil, nil
	}

	validateReq := &BuildValidate{
		Mapping:  m,
		Policies: p,
	}

	data, err := json.Marshal(validateReq)
	if err != nil {
		return nil, err
	}

	msg, err = N.Request("build.validate", data, 2*time.Second)
	if err != nil {
		return nil, err
	}

	var bvr BuildValidateResponse
	err = json.Unmarshal(msg.Data, &bvr)
	if err != nil {
		return nil, err
	}

	return &bvr, nil
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

func (b *BuildValidateResponse) Pass() bool {
	for _, e := range b.Controls {
		if e.Status == "failed" {
			return false
		}
	}

	return true
}
