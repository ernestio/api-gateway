/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import "github.com/ernestio/mapping/validation"

// BuildDetails : for returning build information to the user
type BuildDetails struct {
	ID         string                 `json:"id,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Validation *validation.Validation `json:"validation"`
}
