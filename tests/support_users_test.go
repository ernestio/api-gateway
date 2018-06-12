package main

import (
	"encoding/json"
	"log"

	"github.com/ernestio/api-gateway/models"
	"github.com/nats-io/go-nats"
)

var pw1 = "test1234"
var pw2 = "7bxJZHmyYoEZ9EcjJA85+XXUnZpTcV1Dk4ELpNCfTBeAYougTOcX2+0c/L+jGWuYOLr/UxnmzktpbpQLrBrcaQ=="

var (
	mockUsers = []models.User{
		models.User{
			ID:        1,
			Username:  "test",
			Password:  &pw1,
			MFASecret: "secret",
		},
		models.User{
			ID:       2,
			Username: "test2",
			Password: &pw2,
			Salt:     "yk5c0e8MDkehRv6K83x/5zdo4Jrj/EGBzqTimzJStOo=",
		},
	}
)

func getUserSubscriber(max int) {
	sub, _ := models.N.Subscribe("user.get", func(msg *nats.Msg) {
		var data []byte
		var qu models.User

		if len(msg.Data) > 0 {
			if err := json.Unmarshal(msg.Data, &qu); err != nil {
				log.Println(err)
			}

			for _, user := range mockUsers {
				if user.ID == qu.ID ||
					user.Username == qu.Username {
					data, _ = json.Marshal(user)
					break
				}
			}
			if data != nil {
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

func findUserSubscriber() {
	sub, _ := models.N.Subscribe("user.find", func(msg *nats.Msg) {
		var qu models.User
		var ur []models.User

		if len(msg.Data) == 0 {
			data, _ := json.Marshal(mockUsers)
			if err := models.N.Publish(msg.Reply, data); err != nil {
				log.Println(err)
			}
			return
		}

		if err := json.Unmarshal(msg.Data, &qu); err != nil {
			log.Println(err)
		}

		for _, user := range mockUsers {
			if user.Username == qu.Username || user.ID == qu.ID {
				ur = append(ur, user)
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

func setUserSubscriber() {
	sub, _ := models.N.Subscribe("user.set", func(msg *nats.Msg) {
		var u models.User

		if err := json.Unmarshal(msg.Data, &u); err != nil {
			log.Println(err)
		}

		if u.ID == 0 {
			u.ID = 3
		}

		if u.MFA != nil && *u.MFA {
			u.MFASecret = "secret"
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

func deleteUserSubscriber() {
	sub, _ := models.N.Subscribe("user.del", func(msg *nats.Msg) {
		var u models.User
		if err := json.Unmarshal(msg.Data, &u); err != nil {
			log.Println(err)
		}

		for _, user := range mockUsers {
			if user.ID == u.ID {
				if err := models.N.Publish(msg.Reply, []byte{}); err != nil {
					log.Println(err)
				}
				return
			}
		}

		if err := models.N.Publish(msg.Reply, []byte(`{"_error": "Not found"}`)); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}
