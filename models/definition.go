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
func (d *Definition) MapImport(body []byte) (map[string]interface{}, error) {
	return d.mapDefinition("definition.map.import", body)
}

// MapCreation : calls definition.map.create
func (d *Definition) MapCreation(body []byte) (map[string]interface{}, error) {
	return d.mapDefinition("definition.map.creation", body)
}

// MapDeletion : calls definition.map.deletion
func (d *Definition) MapDeletion(body []byte) (map[string]interface{}, error) {
	return d.mapDefinition("definition.map.deletion", body)
}

// MapCreation : Calls given subject with given body
func (d *Definition) mapDefinition(subject string, body []byte) (map[string]interface{}, error) {
	var m map[string]interface{}

	msg, err := N.Request(subject, body, 1*time.Second)
	if err != nil {
		h.L.Error(err.Error())
		return m, errors.New("Provided yaml is not valid")
	}

	if err := json.Unmarshal(msg.Data, &m); err != nil {
		h.L.Error("Unexpected response from definition.map.creation " + string(msg.Data))
		return m, err
	}

	if m["error"] != nil {
		return m, errors.New(m["error"].(string))
	}

	return m, nil
}
