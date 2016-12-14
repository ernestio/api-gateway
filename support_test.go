/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type handle func(c echo.Context) error

func doRequest(method string, path string, params map[string]string, data []byte, fn handle, ft *jwt.Token) ([]byte, error) {
	var headers map[string]string
	return doRequestHeaders(method, path, params, data, fn, ft, headers)
}

func doRequestHeaders(method string, path string, params map[string]string, data []byte, fn handle, ft *jwt.Token, headers map[string]string) ([]byte, error) {
	e := echo.New()
	req, _ := http.NewRequest(method, path, bytes.NewReader(data))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, echo.NewResponse(rec, e))

	if ft == nil {
		ft = generateTestToken(1, "admin", true)
	}
	c.Set("user", ft)

	for k, v := range params {
		c.SetParamNames(k)
		c.SetParamValues(v)
	}

	c.SetPath(path)
	if err := fn(c); err != nil {
		return []byte(""), err
	}

	resp := rec.Body.Bytes()
	return resp, nil
}

func testsSetup() {
	if err := os.Setenv("JWT_SECRET", "test"); err != nil {
		log.Println(err)
	}
	if err := os.Setenv("NATS_URI", os.Getenv("NATS_URI_TEST")); err != nil {
		log.Println(err)
	}
}
