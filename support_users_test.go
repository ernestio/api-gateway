package main

import (
	"encoding/json"

	"github.com/nats-io/nats"
)

var (
	mockUsers = []User{
		User{
			ID:       1,
			GroupID:  1,
			Username: "test",
			Password: "test",
		},
		User{
			ID:       2,
			GroupID:  2,
			Username: "test2",
			Password: "b3nBt+fHNNSaP2SDeJzNNFfEOiMkqgLh8M7Bajfj2jZZtLp36vAhDMH6i3GXp/EMWTBuBIfQJIA3kgOFqfra0w==",
			Salt:     "psDFaNEE5D9IqCeRrlOmNsRuCKQplicvvXtFhX5S4oE=",
		},
	}
)

func getUserSubscriber(max int) {
	sub, _ := n.Subscribe("user.get", func(msg *nats.Msg) {
		var qu User

		if len(msg.Data) > 0 {
			json.Unmarshal(msg.Data, &qu)

			for _, user := range mockUsers {
				if qu.GroupID != 0 {
					if user.ID == qu.ID && user.GroupID == qu.GroupID ||
						user.Username == qu.Username && user.GroupID == qu.GroupID {
						data, _ := json.Marshal(user)
						n.Publish(msg.Reply, data)
						return
					}
				} else {
					if user.ID == qu.ID || user.Username == qu.Username {
						data, _ := json.Marshal(user)
						n.Publish(msg.Reply, data)
						return
					}
				}
			}
		}

		n.Publish(msg.Reply, []byte(`{"_error":"Not found"}`))
	})
	sub.AutoUnsubscribe(max)
}

func findUserSubscriber() {
	sub, _ := n.Subscribe("user.find", func(msg *nats.Msg) {
		var qu User
		var ur []User

		if len(msg.Data) == 0 {
			data, _ := json.Marshal(mockUsers)
			n.Publish(msg.Reply, data)
			return
		}

		json.Unmarshal(msg.Data, &qu)

		for _, user := range mockUsers {
			if user.Username == qu.Username || user.GroupID == qu.GroupID || user.ID == qu.ID {
				ur = append(ur, user)
			}
		}

		data, _ := json.Marshal(ur)
		n.Publish(msg.Reply, data)
	})
	sub.AutoUnsubscribe(1)
}

func setUserSubscriber() {
	sub, _ := n.Subscribe("user.set", func(msg *nats.Msg) {
		var u User

		json.Unmarshal(msg.Data, &u)
		if u.ID == 0 {
			u.ID = 3
		}

		data, _ := json.Marshal(u)
		n.Publish(msg.Reply, data)
	})
	sub.AutoUnsubscribe(1)
}

func deleteUserSubscriber() {
	sub, _ := n.Subscribe("user.del", func(msg *nats.Msg) {
		var u User
		json.Unmarshal(msg.Data, &u)

		for _, user := range mockUsers {
			if user.ID == u.ID {
				n.Publish(msg.Reply, []byte{})
				return
			}
		}

		n.Publish(msg.Reply, []byte(`{"_error": "Not found"}`))
	})
	sub.AutoUnsubscribe(1)
}
