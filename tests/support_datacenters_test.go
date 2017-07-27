/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"

	"github.com/ernestio/api-gateway/models"
	"github.com/nats-io/nats"
)

var (
	mockDatacenters = []models.Datacenter{
		models.Datacenter{
			ID:   1,
			Name: "test",
		},
		models.Datacenter{
			ID:   2,
			Name: "test2",
		},
	}
)

func getDatacenterSubscriber(max int) {
	sub, _ := models.N.Subscribe("datacenter.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qd := models.Datacenter{}
			if err := json.Unmarshal(msg.Data, &qd); err != nil {
				log.Println(err)
				return
			}

			for _, datacenter := range mockDatacenters {
				if datacenter.ID == qd.ID {
					data, _ := json.Marshal(datacenter)
					if err := models.N.Publish(msg.Reply, data); err != nil {
						log.Println(err)
					}
					return
				}
				data, _ := json.Marshal(datacenter)
				if err := models.N.Publish(msg.Reply, data); err != nil {
					log.Println(err)
				}
				return
			}
		}
		if err := models.N.Publish(msg.Reply, []byte(`{"_error":"Not found"}`)); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(max); err != nil {
		log.Println(err)
	}
}

func findDatacenterSubscriber() {
	sub, _ := models.N.Subscribe("datacenter.find", func(msg *nats.Msg) {
		data, _ := json.Marshal(mockDatacenters)
		if err := models.N.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func createDatacenterSubscriber() {
	sub, _ := models.N.Subscribe("datacenter.set", func(msg *nats.Msg) {
		var d models.Datacenter

		if err := json.Unmarshal(msg.Data, &d); err != nil {
			log.Println(err)
		}
		d.ID = 3
		data, _ := json.Marshal(d)

		if err := models.N.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func deleteDatacenterSubscriber() {
	sub, _ := models.N.Subscribe("datacenter.del", func(msg *nats.Msg) {
		var u models.Datacenter

		if err := json.Unmarshal(msg.Data, &u); err != nil {
			log.Println(err)
		}

		if err := models.N.Publish(msg.Reply, []byte{}); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}
