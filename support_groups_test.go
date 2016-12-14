package main

import (
	"encoding/json"
	"log"

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

func getGroupSubscriber() {
	sub, _ := n.Subscribe("group.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qg := Group{}
			if err := json.Unmarshal(msg.Data, &qg); err != nil {
				log.Println(err)
				return
			}

			for _, group := range mockGroups {
				if group.ID == qg.ID || group.Name == qg.Name {
					data, _ := json.Marshal(group)
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
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func createGroupSubscriber() {
	sub, _ := n.Subscribe("group.set", func(msg *nats.Msg) {
		var g Group

		if err := json.Unmarshal(msg.Data, &g); err != nil {
			log.Println(err)
		}
		g.ID = 3
		data, _ := json.Marshal(g)

		if err := n.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func findGroupSubscriber() {
	sub, _ := n.Subscribe("group.find", func(msg *nats.Msg) {
		var qu Group
		var ur []Group

		if len(msg.Data) == 0 {
			data, _ := json.Marshal(mockGroups)
			if err := n.Publish(msg.Reply, data); err != nil {
				log.Println(err)
			}
			return
		}

		if err := json.Unmarshal(msg.Data, &qu); err != nil {
			log.Println(err)
		}

		for _, group := range mockGroups {
			if group.Name == qu.Name || group.ID == qu.ID {
				ur = append(ur, group)
			}
		}

		data, _ := json.Marshal(ur)
		if err := n.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func setGroupSubscriber() {
	sub, _ := n.Subscribe("group.set", func(msg *nats.Msg) {
		var u Group

		if err := json.Unmarshal(msg.Data, &u); err != nil {
			log.Println(err)
		}
		if u.ID == 0 {
			u.ID = 3
		}

		data, _ := json.Marshal(u)
		if err := n.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func deleteGroupSubscriber() {
	sub, _ := n.Subscribe("group.del", func(msg *nats.Msg) {
		var g Group

		if err := json.Unmarshal(msg.Data, &g); err != nil {
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
