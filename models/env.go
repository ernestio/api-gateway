/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/sirupsen/logrus"
)

// EnvNameSeparator : environment name separator
var EnvNameSeparator = "/"

// Env holds the environment response from service-store
type Env struct {
	ID          uint                   `json:"id"`
	ProjectID   uint                   `json:"project_id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Status      string                 `json:"status"`
	Options     map[string]interface{} `json:"option"`
	Credentials map[string]interface{} `json:"credentials"`
	Roles       []string               `json:"roles,omitempty"`
}

// Validate the env
func (s *Env) Validate() error {
	if s.Name == "" {
		return errors.New("Environment name is empty")
	}

	if s.ProjectID == 0 {
		return errors.New("Environment project is empty")
	}

	if s.Type == "" {
		return errors.New("Environment type is empty")
	}

	return nil
}

// Map : maps a env from a request's body and validates the input
func (s *Env) Map(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	if err := s.Validate(); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Warning("Input is not valid")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// Find : Searches for all envs with filters
func (s *Env) Find(query map[string]interface{}, envs *[]Env) (err error) {
	if err := NewBaseModel(s.getStore()).FindBy(query, envs); err != nil {
		return err
	}
	return nil
}

// FindByName : Searches for all envs with a name equal to the specified
func (s *Env) FindByName(name string) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel(s.getStore()).GetBy(query, s); err != nil {
		return err
	}
	return nil
}

// FindByID : Gets a model by its id
func (s *Env) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	if err := NewBaseModel(s.getStore()).GetBy(query, s); err != nil {
		return err
	}
	return nil
}

// FindAll : Searches for all envs s on the store current user
// has access to
func (s *Env) FindAll(au User, envs *[]Env) (err error) {
	query := make(map[string]interface{})
	if err := NewBaseModel(s.getStore()).FindBy(query, envs); err != nil {
		return err
	}
	return nil
}

// Save : calls env.set with the marshalled
func (s *Env) Save() (err error) {
	if err := NewBaseModel(s.getStore()).Save(s); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a env by its id
func (s *Env) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = s.ID
	if err := NewBaseModel(s.getStore()).Delete(query); err != nil {
		return err
	}
	return nil
}

// DeleteByName : will delete a env by its name
func (s *Env) DeleteByName(name string) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel(s.getStore()).Delete(query); err != nil {
		return err
	}
	return nil
}

// Reset : will reset the env status to errored
func (s *Env) Reset() error {
	var r map[string]interface{}

	s.Status = "errored"
	query := make(map[string]interface{})
	query["id"] = s.ID
	query["status"] = "errored"

	data, err := json.Marshal(query)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}

	resp, err := N.Request("build.set.status", data, time.Second*5)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}

	err = json.Unmarshal(resp.Data, &r)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}

	if r["error"] != nil {
		err = errors.New(r["error"].(string))
		h.L.Error(err.Error())
		return err
	}

	return nil
}

// FindByProjectID : find a envs for the given project id
func (s *Env) FindByProjectID(id int, envs *[]Env) (err error) {
	query := make(map[string]interface{})
	query["project_id"] = id
	if err := NewBaseModel(s.getStore()).FindBy(query, envs); err != nil {
		return err
	}
	return nil
}

/*

// RequestSync : calls service.sync with the given raw message
func (s *Env) RequestSync() error {
	if err := N.Publish(s.getStore()+".sync", []byte(`{"id":"`+s.ID+`"}`)); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

*/

// GetID : ID getter
func (s *Env) GetID() string {
	return s.Name
}

// GetType : Gets the resource type
func (s *Env) GetType() string {
	return "environment"
}

// getStore : Gets the store name
func (s *Env) getStore() string {
	return "environment"
}
