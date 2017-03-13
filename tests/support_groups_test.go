package main

import (
	"encoding/json"
	"log"

	"github.com/ernestio/api-gateway/models"
	"github.com/nats-io/nats"
)

var (
	mockGroups = []models.Group{
		models.Group{
			ID:   1,
			Name: "test",
		},
		models.Group{
			ID:   2,
			Name: "test2",
		},
	}
)

func getGroupSubscriber() {
	sub, _ := models.N.Subscribe("group.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qg := models.Group{}
			if err := json.Unmarshal(msg.Data, &qg); err != nil {
				log.Println(err)
				return
			}

			for _, group := range mockGroups {
				if group.ID == qg.ID || group.Name == qg.Name {
					data, _ := json.Marshal(group)
					if err := models.N.Publish(msg.Reply, data); err != nil {
						log.Println(err)
					}
					return
				}
			}
		}
		if err := models.N.Publish(msg.Reply, []byte(`{"_error":"Not found"}`)); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func createGroupSubscriber() {
	sub, _ := models.N.Subscribe("group.set", func(msg *nats.Msg) {
		var g models.Group

		if err := json.Unmarshal(msg.Data, &g); err != nil {
			log.Println(err)
		}
		g.ID = 3
		data, _ := json.Marshal(g)

		if err := models.N.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func findGroupSubscriber() {
	sub, _ := models.N.Subscribe("group.find", func(msg *nats.Msg) {
		var qu models.Group
		var ur []models.Group

		if len(msg.Data) == 0 {
			data, _ := json.Marshal(mockGroups)
			if err := models.N.Publish(msg.Reply, data); err != nil {
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
		if err := models.N.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func setGroupSubscriber() {
	sub, _ := models.N.Subscribe("group.set", func(msg *nats.Msg) {
		var u models.Group

		if err := json.Unmarshal(msg.Data, &u); err != nil {
			log.Println(err)
		}
		if u.ID == 0 {
			u.ID = 3
		}

		data, _ := json.Marshal(u)
		if err := models.N.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func deleteGroupSubscriber() {
	sub, _ := models.N.Subscribe("group.del", func(msg *nats.Msg) {
		var g models.Group

		if err := json.Unmarshal(msg.Data, &g); err != nil {
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
