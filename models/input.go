package models

// ServiceInput : service received by the endpoint
type ServiceInput struct {
	Datacenter string `json:"datacenter"`
	Name       string `json:"name"`
}
