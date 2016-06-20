/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

type handle func(c echo.Context) error

func doRequest(method string, path string, params map[string]string, data []byte, fn handle, ft *jwt.Token) ([]byte, error) {
	e := echo.New()
	req, _ := http.NewRequest(method, path, bytes.NewReader(data))
	rec := httptest.NewRecorder()
	c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

	if ft == nil {
		ft = jwt.New(jwt.SigningMethodHS256)
		ft.Claims["username"] = "admin"
		ft.Claims["admin"] = true
		ft.Claims["group_id"] = 2.0
	}
	c.Set("user", ft)

	for k, v := range params {
		c.SetParamNames(k)
		c.SetParamValues(v)
	}

	c.SetPath(path)
	if err := fn(c); err != nil {
		return []byte(""), err
	} else {
		resp := rec.Body.Bytes()
		return resp, nil
	}
}
