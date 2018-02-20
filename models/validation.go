package models

// BuildValidate describes a request to the build validate service.
type BuildValidate struct {
	Mapping  *Mapping `json:"mapping"`
	Policies []Policy `json:"policies"`
}

// BuildValidateReponse describes a response from the build validate service.
type BuildValidateResponse struct {
	Version    string     `json:"version"`
	Controls   []Control  `json:"controls"`
	Profiles   []Profile  `json:"profiles"`
	Statistics Statistics `json:"statistics"`
}

// Profile describes the policy document and its test results
type Profile struct {
	Supports []string         `json:"supports"`
	Title    string           `json:"title"`
	Name     string           `json:"name"`
	Controls []ControlDetails `json:"controls"`
}

// ControlDetails describes additional information about a control
type ControlDetails struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"desc"`
	Impact         float32   `json:"impact"`
	References     []string  `json:"refs"`
	Tags           []string  `json:"tags"`
	Code           string    `json:"code"`
	Results        []Control `json:"results"`
	Groups         []Group   `json:"groups"`
	Attributes     []string  `json:"attribues"`
	SHA256         string    `json:"sha256"`
	SourceLocation struct {
		Reference string `json:"ref"`
		Line      int    `json:"line"`
	} `json:"source_location"`
}

// Control describes an individual test within a build validation.
type Control struct {
	ID        string  `json:"id"`
	ProfileID string  `json:"profile_id"`
	Status    string  `json:"status"`
	CodeDesc  string  `json:"code_desc"`
	RunTime   float64 `json:"run_time"`
	StartTime string  `json:"start_time"`
	Message   string  `json:"message"`
}

type Group struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Controls []string `json:"controls"`
}

// Statistics describes stats for the build validate service.
type Statistics struct {
	Duration float64 `json:"duration"`
}

// Passed : returns true if validation rules passed
func (b *BuildValidateResponse) Passed() bool {
	for _, e := range b.Controls {
		if e.Status == "failed" {
			return false
		}
	}

	return true
}
