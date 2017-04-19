/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package views

import (
	"encoding/json"
	"time"

	"github.com/ernestio/api-gateway/models"
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

// RenderUsageReport : TODO
func RenderUsageReport(reportables []models.Usage) ([]byte, error) {
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

func mapUsage(u models.Usage) UsageRender {
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
