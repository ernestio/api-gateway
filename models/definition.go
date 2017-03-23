package models

import (
	"encoding/json"
	"errors"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
)

// Definition : definition interface
type Definition struct {
}

// MapImport : calls definition.map.import
func (d *Definition) MapImport(body []byte) ([]byte, error) {
	return d.mapDefinition("definition.map.import", body)
}

// MapCreation : calls definition.map.create
func (d *Definition) MapCreation(body []byte) ([]byte, error) {
	return d.mapDefinition("definition.map.creation", body)
}

// MapDeletion : calls definition.map.deletion
func (d *Definition) MapDeletion(previous, serviceType string) ([]byte, error) {
	query := []byte(`{"previous_id":"` + previous + `","datacenter":{"type":"` + serviceType + `"}}`)
	msg, err := N.Request("definition.map.deletion", query, 1*time.Second)
	if err != nil {
		h.L.Error(err.Error())
		return []byte(""), errors.New("Couldn't map the service")
	}
	return msg.Data, nil
}

// MapCreation : Calls given subject with given body
func (d *Definition) mapDefinition(subject string, body []byte) ([]byte, error) {
	msg, err := N.Request(subject, body, 1*time.Second)
	if err != nil {
		h.L.Error(err.Error())
		return body, errors.New("Provided yaml is not valid")
	}

	var s struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(msg.Data, &s); err != nil {
		h.L.Error("Unexpected response from definition.map.creation " + string(msg.Data))
		return body, err
	}
	if s.Error != "" {
		return body, errors.New(s.Error)
	}

	return msg.Data, nil
}
