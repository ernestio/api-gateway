/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats"
)

var (
	mockDatacenters = []Datacenter{
		Datacenter{
			ID:      1,
			Name:    "test",
			GroupID: 1,
		},
		Datacenter{
			ID:      2,
			Name:    "test2",
			GroupID: 2,
		},
	}
)

func getDatacenterSubscriber(max int) {
	sub, _ := n.Subscribe("datacenter.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qd := Datacenter{}
			if err := json.Unmarshal(msg.Data, &qd); err != nil {
				log.Println(err)
				return
			}

			for _, datacenter := range mockDatacenters {
				if qd.GroupID != 0 && datacenter.GroupID == qd.GroupID && datacenter.ID == qd.ID {
					data, _ := json.Marshal(datacenter)
					if err := n.Publish(msg.Reply, data); err != nil {
						log.Println(err)
					}
					return
				} else if qd.GroupID == 0 && datacenter.ID == qd.ID {
					data, _ := json.Marshal(datacenter)
					if err := n.Publish(msg.Reply, data); err != nil {
						log.Println(err)
					}
					return
				}
			}
		}
		if err := n.Publish(msg.Reply, []byte(`{"_error":"Not found"}`)); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(max); err != nil {
		log.Println(err)
	}
}

func findDatacenterSubscriber() {
	sub, _ := n.Subscribe("datacenter.find", func(msg *nats.Msg) {
		data, _ := json.Marshal(mockDatacenters)
		if err := n.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func createDatacenterSubscriber() {
	sub, _ := n.Subscribe("datacenter.set", func(msg *nats.Msg) {
		var d Datacenter

		if err := json.Unmarshal(msg.Data, &d); err != nil {
			log.Println(err)
		}
		d.ID = 3
		data, _ := json.Marshal(d)

		if err := n.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func deleteDatacenterSubscriber() {
	sub, _ := n.Subscribe("datacenter.del", func(msg *nats.Msg) {
		var u Datacenter

		if err := json.Unmarshal(msg.Data, &u); err != nil {
			log.Println(err)
		}

		if err := n.Publish(msg.Reply, []byte{}); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}
