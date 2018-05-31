/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"log"
	"regexp"
)

func IsAlphaNumeric(s string) bool {
	r, err := regexp.Compile(`[^a-zA-Z0-9@._\-//]+`)
	if err != nil {
		log.Println(err)
		return false
	}

	return !r.MatchString(s)
}
