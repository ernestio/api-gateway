/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
)

type handle func(c echo.Context) error

func testsSetup() {
	if err := os.Setenv("JWT_SECRET", "test"); err != nil {
		log.Println(err)
	}
	if err := os.Setenv("NATS_URI", os.Getenv("NATS_URI_TEST")); err != nil {
		log.Println(err)
	}
}
