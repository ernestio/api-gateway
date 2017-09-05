package models

// ServiceInput : service received by the endpoint
type ServiceInput struct {
	Datacenter  string                 `json:"project"`
	Credentials map[string]interface{} `json:"credentials"`
	Filters     []string               `json:"import_filters"`
	Name        string                 `json:"name"`
}
