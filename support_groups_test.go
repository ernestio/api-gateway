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
			ID:   2,
			Name: "test2",
		},
	}
)

func getGroupSubcriber() {
	n.Subscribe("group.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qg := Group{}
			json.Unmarshal(msg.Data, &qg)

			for _, group := range mockGroups {
				data, _ := json.Marshal(group)
				n.Publish(msg.Reply, data)
				return
			}
		}
		n.Publish(msg.Reply, []byte(`{"error":"not found"}`))
	})
}

func createGroupSubcriber() {
	n.Subscribe("group.set", func(msg *nats.Msg) {
		var g Group

		json.Unmarshal(msg.Data, &g)
		g.ID = 3
		data, _ := json.Marshal(g)

		n.Publish(msg.Reply, data)
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
		var g Group

		json.Unmarshal(msg.Data, &g)

		n.Publish(msg.Reply, []byte{})
	})
}
