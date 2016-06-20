package main

import (
	"encoding/json"

	"github.com/nats-io/nats"
)

var (
	mockGroups = []Group{
		Group{
			ID:   1,
			Name: "test",
		},
		Group{
			ID:   1,
			Name: "test",
		},
	}
)

func getGroupSubcriber() {
	n.Subscribe("group.get", func(msg *nats.Msg) {
		var qu Group

		if len(msg.Data) > 0 {
			json.Unmarshal(msg.Data, &qu)

			for _, group := range mockGroups {
				if group.ID == qu.ID || group.Name == qu.Name {
					data, _ := json.Marshal(group)
					n.Publish(msg.Reply, data)
					return
				}
			}
		}

		n.Publish(msg.Reply, []byte(`{"error":"not found"}`))
	})
}

func findGroupSubcriber() {
	n.Subscribe("group.find", func(msg *nats.Msg) {
		var qu Group
		var ur []Group

		if len(msg.Data) == 0 {
			data, _ := json.Marshal(mockGroups)
			n.Publish(msg.Reply, data)
			return
		}

		json.Unmarshal(msg.Data, &qu)

		for _, group := range mockGroups {
			if group.Name == qu.Name || group.ID == qu.ID {
				ur = append(ur, group)
			}
		}

		data, _ := json.Marshal(ur)
		n.Publish(msg.Reply, data)
	})
}

func setGroupSubcriber() {
	n.Subscribe("group.set", func(msg *nats.Msg) {
		var u Group

		json.Unmarshal(msg.Data, &u)
		if u.ID == 0 {
			u.ID = 3
		}

		data, _ := json.Marshal(u)
		n.Publish(msg.Reply, data)
	})
}

func deleteGroupSubcriber() {
	n.Subscribe("group.del", func(msg *nats.Msg) {
		var u Datacenter

		json.Unmarshal(msg.Data, &u)

		n.Publish(msg.Reply, []byte{})
	})
}
