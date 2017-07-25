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

// Service holds the service response from service-store
type Service struct {
	ID           string      `json:"id"`
	GroupID      int         `json:"group_id"`
	UserID       int         `json:"user_id"`
	UserName     string      `json:"user_name,omitempty"`
	DatacenterID int         `json:"datacenter_id"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Version      time.Time   `json:"version"`
	Options      string      `json:"options"`
	Status       string      `json:"status"`
	Endpoint     string      `json:"endpoint"`
	Definition   interface{} `json:"definition"`
	Maped        string      `json:"mapping"`
	Sync         bool        `json:"sync"`
	SyncType     string      `json:"sync_type"`
	SyncInterval int         `json:"sync_interval"`
}

// Validate the service
func (s *Service) Validate() error {
	if s.Name == "" {
		return errors.New("Service name is empty")
	}

	if s.DatacenterID == 0 {
		return errors.New("Service group is empty")
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

// FindByName : Searches for all services with a name equal to the specified
func (s *Service) FindByName(name string, service *Service) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel("service").GetBy(query, service); err != nil {
		return err
	}
	return nil
}

// FindByGroupID : Searches for all services on the store current user
// has access to with the specified group id
func (s *Service) FindByGroupID(id int, services *[]Service) (err error) {
	query := make(map[string]interface{})
	query["group_id"] = id

	return NewBaseModel("service").FindBy(query, services)
}

// FindByNameAndGroupID : Searches for all services with a name equal to the specified
func (s *Service) FindByNameAndGroupID(name string, id int, service *[]Service) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	query["group_id"] = id

	return NewBaseModel("service").FindBy(query, service)
}

// GetByNameAndGroupID : Searches for all services with a name equal to the specified
func (s *Service) GetByNameAndGroupID(name string, group int) (service *Service, err error) {
	var services []Service

	if err = s.FindByNameAndGroupID(name, group, &services); err != nil {
		return service, h.ErrGatewayTimeout
	}

	if len(services) == 0 {
		return nil, nil
	}

	return &services[0], nil
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

// FindAll : Searches for all groups on the store current user
// has access to
func (s *Service) FindAll(au User, services *[]Service) (err error) {
	query := make(map[string]interface{})
	query["group_id"] = au.GroupID
	if err := NewBaseModel("service").FindBy(query, services); err != nil {
		return err
	}
	return nil
}

// Save : calls service.set with the marshalled current group
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
func (s *Service) Reset() (err error) {
	s.Status = "errored"
	query := make(map[string]interface{})
	query["id"] = s.ID
	query["status"] = "errored"

	err = NewBaseModel("service").Set(query)

	return err
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
func (s *Service) RequestCreation(raw []byte) error {
	if err := N.Publish("service.create", raw); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestImport : calls service.import with the given raw message
func (s *Service) RequestImport(raw []byte) error {
	if err := N.Publish("service.import", raw); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

// RequestDeletion : calls service.delete with the given raw message
func (s *Service) RequestDeletion(raw []byte) error {
	if err := N.Publish("service.delete", raw); err != nil {
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
