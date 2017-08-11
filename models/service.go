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

// Service holds the service response from service-store
type Service struct {
	ID           string                 `json:"id"`
	UserID       int                    `json:"user_id"`
	UserName     string                 `json:"user_name,omitempty"`
	Project      string                 `json:"project,omitempty"`
	Provider     string                 `json:"provider,omitempty"`
	DatacenterID int                    `json:"datacenter_id"`
	ProjectInfo  *json.RawMessage       `json:"datacenter_info"`
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

// Validate the service
func (s *Service) Validate() error {
	if s.Name == "" {
		return errors.New("Service name is empty")
	}

	if s.DatacenterID == 0 {
		return errors.New("Service datacenter is empty")
	}

	if s.Type == "" {
		return errors.New("Service type is empty")
	}

	return nil
}

// Map : maps a service from a request's body and validates the input
func (s *Service) Map(data []byte) error {
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

// Find : Searches for all services with filters
func (s *Service) Find(query map[string]interface{}, services *[]Service) (err error) {
	if err := NewBaseModel("service").FindBy(query, services); err != nil {
		return err
	}
	return nil
}

// FindLastByName : Searches for all services with a name equal to the specified
func (s *Service) FindLastByName(name string) (service Service, err error) {
	var ss []Service
	query := make(map[string]interface{})
	query["name"] = name
	err = s.Find(query, &ss)
	if len(ss) > 0 {
		service = ss[0]
	}

	return
}

// FindByName : Searches for all services with a name equal to the specified
func (s *Service) FindByName(name string, service *Service) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel("service").GetBy(query, service); err != nil {
		return err
	}
	return nil
}

// FindByID : Gets a model by its id
func (s *Service) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	if err := NewBaseModel("service").GetBy(query, s); err != nil {
		return err
	}
	return nil
}

// FindAll : Searches for all services s on the store current user
// has access to
func (s *Service) FindAll(au User, services *[]Service) (err error) {
	query := make(map[string]interface{})
	if err := NewBaseModel("service").FindBy(query, services); err != nil {
		return err
	}
	return nil
}

// Save : calls service.set with the marshalled
func (s *Service) Save() (err error) {
	if err := NewBaseModel("service").Save(s); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a service by its id
func (s *Service) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = s.ID
	if err := NewBaseModel("service").Delete(query); err != nil {
		return err
	}
	return nil
}

// DeleteByName : will delete a service by its name
func (s *Service) DeleteByName(name string) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel("service").Delete(query); err != nil {
		return err
	}
	return nil
}

// Mapping : will get a service mapping
func (s *Service) Mapping() (*graph.Graph, error) {
	var m map[string]interface{}

	query := make(map[string]interface{})
	query["id"] = s.ID

	err := NewBaseModel("service").CallStoreBy("get.mapping", query, &m)
	if err != nil {
		return nil, err
	}

	g := graph.New()
	err = g.Load(m)

	return g, err
}

// Reset : will reset the service status to errored
func (s *Service) Reset() error {
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

// FindByDatacenterID : find a services for the given datacenter id
func (s *Service) FindByDatacenterID(id int, services *[]Service) (err error) {
	query := make(map[string]interface{})
	query["datacenter_id"] = id
	if err := NewBaseModel("service").FindBy(query, services); err != nil {
		return err
	}
	return nil
}

// RequestCreation : calls service.create with the given raw message
func (s *Service) RequestCreation(mapping map[string]interface{}) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish("service.create", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestImport : calls service.import with the given raw message
func (s *Service) RequestImport(mapping map[string]interface{}) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish("service.import", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestDeletion : calls service.delete with the given raw message
func (s *Service) RequestDeletion(mapping map[string]interface{}) error {
	data, err := json.Marshal(mapping)
	if err != nil {
		h.L.Error(err.Error())
		return err
	}
	if err := N.Publish("service.delete", data); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestSync : calls service.sync with the given raw message
func (s *Service) RequestSync() error {
	if err := N.Publish("service.sync", []byte(`{"id":"`+s.ID+`"}`)); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// GetID : ID getter
func (s *Service) GetID() string {
	return s.Name
}

// GetType : Gets the resource type
func (s *Service) GetType() string {
	return "environment"
}
