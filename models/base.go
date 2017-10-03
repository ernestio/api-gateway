/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/r3labs/akira"
	"github.com/sirupsen/logrus"
)

// N : Nats connection
var N akira.Connector

// BaseModel : Abstraction layer to interact with data stores
type BaseModel struct {
	Type string
}

// NewBaseModel : Constructor
func NewBaseModel(t string) *BaseModel {
	return &BaseModel{Type: t}
}

// CallStoreBy : ...
func (b *BaseModel) CallStoreBy(verb string, query map[string]interface{}, o interface{}) (err error) {
	var res []byte

	err = b.CallStoreByRaw(verb, query, &res)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(res, &o); err != nil {
		return errors.New(`"Specified ` + b.Type + ` does not exist"`)
	}

	return nil
}

// CallStoreByRaw : ...
func (b *BaseModel) CallStoreByRaw(verb string, query map[string]interface{}, res *[]byte) (err error) {
	var req []byte
	var rm map[string]interface{}

	if len(query) > 0 {
		if req, err = json.Marshal(query); err != nil {
			return err
		}
	}

	if *res, err = b.Query(b.Type+"."+verb, string(req)); err != nil {
		return err
	}

	json.Unmarshal(*res, rm)
	if rm["error"] != nil {
		return errors.New(rm["error"].(string))
	}

	return nil
}

// GetBy : interface to call component.get on the specific store
func (b *BaseModel) GetBy(query map[string]interface{}, o interface{}) (err error) {
	return b.CallStoreBy("get", query, o)
}

// FindBy : interface to call component.find on the specific store
func (b *BaseModel) FindBy(query map[string]interface{}, o interface{}) (err error) {
	return b.CallStoreBy("find", query, o)
}

// Save : interface to call component.set on the specific store
func (b *BaseModel) Save(o interface{}) (err error) {
	var res []byte

	data, err := json.Marshal(o)
	if err != nil {
		return NewError(InvalidInputCode, "Can't marshal provided component to json")
	}

	if res, err = b.Query(b.Type+".set", string(data)); err != nil {
		return err
	}
	if err := json.Unmarshal(res, &o); err != nil {
		msg := "An internal error occurred saving component " + b.Type
		h.L.WithFields(logrus.Fields{
			"response":      string(res),
			"error_message": err.Error(),
		}).Error(msg)
		return NewError(InternalCode, msg)
	}

	return nil
}

// Delete : interface to call component.del on the specific store
func (b *BaseModel) Delete(query map[string]interface{}) (err error) {
	var res []byte
	var req []byte
	if len(query) > 0 {
		if req, err = json.Marshal(query); err != nil {
			return err
		}
	}
	if res, err = b.Query(b.Type+".del", string(req)); err != nil {
		return err
	}
	if strings.Contains(string(res), `"error"`) {
		return errors.New(`"Specified ` + b.Type + ` does not exist"`)
	}

	return nil
}

// Query : Allows a free query by subject
func (b *BaseModel) Query(subject, query string) ([]byte, error) {
	var res []byte
	msg, err := N.Request(subject, []byte(query), 5*time.Second)
	if err != nil {
		eMsg := "An internal error happened trying to reach the data store"
		h.L.Error(eMsg)
		return res, NewError(TimeoutCode, eMsg)
	}

	if re := h.ResponseErr(msg); re != nil {
		return res, re.HTTPError
	}

	return msg.Data, nil
}

// Set : interface to call component.set on the specific store
func (b *BaseModel) Set(query map[string]interface{}) (err error) {
	var req []byte
	if len(query) > 0 {
		if req, err = json.Marshal(query); err != nil {
			return err
		}
	}

	if _, err = b.Query(b.Type+".set", string(req)); err != nil {
		return err
	}

	return nil
}
