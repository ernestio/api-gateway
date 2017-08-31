package models

// ServiceInput : service received by the endpoint
type ServiceInput struct {
	Datacenter  string                 `json:"project"`
	Credentials map[string]interface{} `json:"credentials"`
	Name        string                 `json:"name"`
}
