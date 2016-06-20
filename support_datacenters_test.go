/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"

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

func getDatacenterSubcriber() {
	n.Subscribe("datacenter.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qd := Datacenter{}
			json.Unmarshal(msg.Data, &qd)

			for _, datacenter := range mockDatacenters {
				if qd.GroupID != 0 && datacenter.GroupID == qd.GroupID && datacenter.Name == qd.Name {
					data, _ := json.Marshal(datacenter)
					n.Publish(msg.Reply, data)
					return
				} else if qd.GroupID == 0 && datacenter.Name == qd.Name {
					data, _ := json.Marshal(datacenter)
					n.Publish(msg.Reply, data)
					return
				}
			}
		}
		n.Publish(msg.Reply, []byte(`{"error":"not found"}`))
	})
}

func findDatacenterSubcriber() {
	n.Subscribe("datacenter.find", func(msg *nats.Msg) {
		data, _ := json.Marshal(mockDatacenters)
		n.Publish(msg.Reply, data)
	})
}

func createDatacenterSubcriber() {
	n.Subscribe("datacenter.set", func(msg *nats.Msg) {
		var d Datacenter

		json.Unmarshal(msg.Data, &d)
		d.ID = 3
		data, _ := json.Marshal(d)

		n.Publish(msg.Reply, data)
	})
}

func deleteDatacenterSubcriber() {
	n.Subscribe("datacenter.del", func(msg *nats.Msg) {
		var u Datacenter

		json.Unmarshal(msg.Data, &u)

		n.Publish(msg.Reply, []byte{})
	})
}
