/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package config

import (
	"log"

	"github.com/ernestio/api-gateway/controllers"
	"github.com/ernestio/api-gateway/models"
	ecc "github.com/ernestio/ernest-config-client"
)

// Setup : Set up api-gateway based on its environment requirements
func Setup() {
	var err error

	log.Println("Getting configuration parameters")
	c := models.Config{}
	natsuri := c.GetNatsURI()
	models.N = ecc.NewConfig(natsuri).Nats()
	if controllers.Secret, err = c.GetJWTToken(); err != nil {
		panic(err.Error())
	}
}
