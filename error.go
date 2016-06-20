/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"strings"

	"github.com/labstack/echo"
	"github.com/nats-io/nats"
)

// ResponseError is
type ResponseError struct {
	Error     string          `json:"error"`
	HTTPError *echo.HTTPError `json:"-"`
}

func responseErr(msg *nats.Msg) *ResponseError {
	var e ResponseError

	err := json.Unmarshal(msg.Data, &e)
	if err != nil || e.Error == "" {
		return nil
	}

	if strings.Contains(e.Error, "not found") {
		e.HTTPError = ErrNotFound
	}

	if strings.Contains(e.Error, "unexpected") {
		e.HTTPError = ErrInternal
	}

	return &e
}
