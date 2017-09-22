/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/r3labs/graph"
	"github.com/sirupsen/logrus"
)

// Build : holds the build response from service store
type Build struct {
	ID            string                 `json:"id"`
	EnvironmentID int                    `json:"environment_id"`
	UserID        int                    `json:"user_id"`
	Username      string                 `json:"user_name"`
	Type          string                 `json:"type"`
	Status        string                 `json:"status"`
	Definition    string                 `json:"definition"`
	Mapping       map[string]interface{} `json:"mapping"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// Validate the env
func (b *Build) Validate() error {
	if b.EnvironmentID == 0 {
		return errors.New("Build environment id is empty")
	}

	if b.Type == "" {
		return errors.New("Build type is empty")
	}

	return nil
}

// Map : maps a build from a request's body and validates the input
func (b *Build) Map(data []byte) error {
	if err := json.Unmarshal(data, b); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	if err := b.Validate(); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Warning("Input is not valid")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// Find : Searches for all builds with filters
func (b *Build) Find(query map[string]interface{}, builds *[]Build) (err error) {
	if err := NewBaseModel(b.getStore()).FindBy(query, builds); err != nil {
		return err
	}
	return nil
}

// FindLastByName : Searches for all environments with a name equal to the specified
func (b *Build) FindLastByName(env string) (err error) {
	var e Env

	err = e.FindByName(env)
	if err != nil {
		return
	}

	query := make(map[string]interface{})
	query["environment_id"] = e.ID
	if err = NewBaseModel(b.getStore()).GetBy(query, b); err != nil {
		return
	}

	return
}

// FindByName : Searches for all builds with by env name
func (b *Build) FindByName(env string) (builds []Build, err error) {
	var e Env

	err = e.FindByName(env)
	if err != nil {
		return
	}

	query := make(map[string]interface{})
	query["environment_id"] = e.ID
	if err = NewBaseModel(b.getStore()).FindBy(query, &builds); err != nil {
		return
	}

	return
}

// FindByID : Gets a model by its id
func (b *Build) FindByID(id string) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	if err := NewBaseModel(b.getStore()).GetBy(query, b); err != nil {
		return err
	}
	return nil
}

// FindAll : Searches for all builds on the store current user
// has access to
func (b *Build) FindAll(au User, builds *[]Build) (err error) {
	query := make(map[string]interface{})
	if err := NewBaseModel(b.getStore()).FindBy(query, builds); err != nil {
		return err
	}
	return nil
}

// Save : calls build.set with the marshalled
func (b *Build) Save() (err error) {
	if err := NewBaseModel(b.getStore()).Save(b); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a builds by its id
func (b *Build) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = b.ID
	if err := NewBaseModel(b.getStore()).Delete(query); err != nil {
		return err
	}
	return nil
}

// GetMapping : will get a builds mapping
func (b *Build) GetMapping() (*graph.Graph, error) {
	var m Mapping
	query := make(map[string]interface{})
	query["id"] = b.ID

	err := NewBaseModel(b.getStore()).CallStoreBy("get.mapping", query, &m)
	if err != nil {
		return nil, err
	}

	g := graph.New()
	err = g.Load(m)

	return g, err
}

// Reset : will reset the builds status to errored
func (b *Build) Reset() error {
	var r map[string]interface{}

	b.Status = "errored"
	query := make(map[string]interface{})
	query["id"] = b.ID
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

// FindByEnvironmentName : find a builds for the given environment name
func (b *Build) FindByEnvironmentName(name string, builds *[]Build) (err error) {
	var e Env

	err = e.FindByName(name)
	if err != nil {
		return err
	}

	query := make(map[string]interface{})
	query["environment_id"] = e.ID
	if err := NewBaseModel(b.getStore()).FindBy(query, builds); err != nil {
		return err
	}

	return nil
}

// RequestCreation : calls env.create with the given raw message
func (b *Build) RequestCreation(mapping *Mapping) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish(b.getStore()+".create", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestImport : calls build.import with the given raw message
func (b *Build) RequestImport(mapping *Mapping) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish(b.getStore()+".import", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestDeletion : calls builddelete with the given raw message
func (b *Build) RequestDeletion(mapping *Mapping) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish(b.getStore()+".delete", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

/*

// RequestSync : calls service.sync with the given raw message
func (b *Build) RequestSync() error {
	if err := N.Publish(s.getStore()+".sync", []byte(`{"id":"`+s.ID+`"}`)); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

*/

// GetID : ID getter
func (b *Build) GetID() string {
	return b.ID
}

// GetType : Gets the resource type
func (b *Build) GetType() string {
	return "build"
}

// getStore : Gets the store name
func (b *Build) getStore() string {
	return "build"
}
