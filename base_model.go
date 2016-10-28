/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Group holds the group response from group-store
type BaseModel struct {
	Type string
}

func NewBaseModel(t string) *BaseModel {
	return &BaseModel{Type: t}
}

func (b *BaseModel) callStoreBy(verb string, query map[string]interface{}, o interface{}) (err error) {
	var res []byte
	var req []byte
	if len(query) > 0 {
		if req, err = json.Marshal(query); err != nil {
			return err
		}
	}
	if res, err = b.Query(b.Type+"."+verb, string(req)); err != nil {
		return err
	}
	if strings.Contains(string(res), `"error"`) {
		return errors.New(`"Specified ` + b.Type + ` does not exist"`)
	}
	if err = json.Unmarshal(res, &o); err != nil {
		return errors.New(`"Specified ` + b.Type + ` does not exist"`)
	}

	return nil
}

// GetBy : interface to call component.get on the specific store
func (b *BaseModel) GetBy(query map[string]interface{}, o interface{}) (err error) {
	return b.callStoreBy("get", query, o)
}

// FindBy : interface to call component.find on the specific store
func (b *BaseModel) FindBy(query map[string]interface{}, o interface{}) (err error) {
	return b.callStoreBy("find", query, o)
}

// Save : interface to call component.set on the specific store
func (b *BaseModel) Save(o interface{}) (err error) {
	var res []byte

	data, err := json.Marshal(o)
	if err != nil {
		return ErrBadReqBody
	}

	if res, err = b.Query(b.Type+".set", string(data)); err != nil {
		return err
	}
	if err := json.Unmarshal(res, &o); err != nil {
		println(string(res))
		println(err.Error())
		return ErrInternal
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

func (b *BaseModel) Query(subject, query string) ([]byte, error) {
	var res []byte
	msg, err := n.Request(subject, []byte(query), 5*time.Second)
	if err != nil {
		return res, ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
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
