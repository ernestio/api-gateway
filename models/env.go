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
	graph "gopkg.in/r3labs/graph.v2"
)

// EnvNameSeparator : environment name separator
var EnvNameSeparator = "/"

// Env holds the environment response from service-store
type Env struct {
	ID           string                 `json:"id"`
	UserID       int                    `json:"user_id"`
	UserName     string                 `json:"user_name,omitempty"`
	Project      string                 `json:"project,omitempty"`
	Provider     string                 `json:"provider,omitempty"`
	DatacenterID int                    `json:"datacenter_id"`
	ProjectInfo  *json.RawMessage       `json:"credentials,omitempty"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Version      time.Time              `json:"version"`
	Options      map[string]interface{} `json:"options"`
	Status       string                 `json:"status"`
	Endpoint     string                 `json:"endpoint"`
	Definition   interface{}            `json:"definition"`
	Mapped       map[string]interface{} `json:"mapping"`
	Sync         bool                   `json:"sync"`
	SyncType     string                 `json:"sync_type"`
	SyncInterval int                    `json:"sync_interval"`
	Roles        []string               `json:"roles,omitempty"`
}

// Validate the env
func (s *Env) Validate() error {
	if s.Name == "" {
		return errors.New("Environment name is empty")
	}

	if s.DatacenterID == 0 {
		return errors.New("Environment datacenter is empty")
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

// FindLastByName : Searches for all environments with a name equal to the specified
func (s *Env) FindLastByName(name string) (env Env, err error) {
	var ss []Env
	query := make(map[string]interface{})
	query["name"] = name
	err = s.Find(query, &ss)
	if len(ss) > 0 {
		env = ss[0]
	}

	return
}

// FindByName : Searches for all envs with a name equal to the specified
func (s *Env) FindByName(name string, env *Env) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel(s.getStore()).GetBy(query, env); err != nil {
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

// Mapping : will get a env mapping
func (s *Env) Mapping() (*graph.Graph, error) {
	var m map[string]interface{}

	query := make(map[string]interface{})
	query["id"] = s.ID

	err := NewBaseModel(s.getStore()).CallStoreBy("get.mapping", query, &m)
	if err != nil {
		return nil, err
	}

	g := graph.New()
	err = g.Load(m)

	return g, err
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
	query["datacenter_id"] = id
	if err := NewBaseModel(s.getStore()).FindBy(query, envs); err != nil {
		return err
	}
	return nil
}

// RequestCreation : calls env.create with the given raw message
func (s *Env) RequestCreation(mapping map[string]interface{}) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish(s.getStore()+".create", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestImport : calls service.import with the given raw message
func (s *Env) RequestImport(mapping map[string]interface{}) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish(s.getStore()+".import", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestDeletion : calls service.delete with the given raw message
func (s *Env) RequestDeletion(mapping map[string]interface{}) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish(s.getStore()+".delete", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestSync : calls service.sync with the given raw message
func (s *Env) RequestSync() error {
	if err := N.Publish(s.getStore()+".sync", []byte(`{"id":"`+s.ID+`"}`)); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

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
	return "service"
}
