package models

import "encoding/json"

// ServiceInput : service received by the endpoint
type ServiceInput struct {
	Datacenter  string           `json:"datacenter"`
	ProjectInfo *json.RawMessage `json:"credentials"`
	Name        string           `json:"name"`
}
