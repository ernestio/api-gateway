package main

import (
	"encoding/json"
	"time"
)

// UsageRender : usage struct to be rendered
type UsageRender struct {
	ID      uint       `json:"id" gorm:"primary_key"`
	Service string     `json:"service"`
	Name    string     `json:"name"`
	Type    string     `json:"type"`
	From    *time.Time `json:"from"`
	To      *time.Time `json:"to"`
}

func renderUsageReport(reportables []Usage) ([]byte, error) {
	result := make(map[string][]UsageRender)

	for _, r := range reportables {
		key := r.Service + "-" + r.Name + "-" + r.Type
		if _, ok := result[key]; !ok {
			result[key] = make([]UsageRender, 0)
		}
		result[key] = append(result[key], mapUsage(r))
	}

	return json.Marshal(result)
}

func mapUsage(u Usage) UsageRender {
	var from, to time.Time
	if u.From != 0 {
		from = time.Unix(u.From, 0)
	}
	if u.To != 0 {
		to = time.Unix(u.To, 0)
	}

	return UsageRender{
		ID:      u.ID,
		Service: u.Service,
		Name:    u.Name,
		Type:    u.Type,
		From:    &from,
		To:      &to,
	}
}
