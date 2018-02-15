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
	Statistics Statistics `json:"statistics"`
}

// Control describes an individual test within a build validation.
type Control struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	Status    string `json:"status"`
	CodeDesc  string `json:"code_desc"`
	Message   string `json:"message"`
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
