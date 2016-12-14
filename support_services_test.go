/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats"
)

var (
	mockServices = []Service{
		Service{
			ID:           "1",
			Name:         "test",
			GroupID:      1,
			DatacenterID: 1,
			Version:      time.Now(),
		},
		Service{
			ID:           "3",
			Name:         "test",
			GroupID:      1,
			DatacenterID: 1,
			Version:      time.Now(),
		},
		Service{
			ID:           "2",
			Name:         "test2",
			GroupID:      2,
			DatacenterID: 3,
			Version:      time.Now(),
		},
	}
)

func getServiceSubscriber() {
	_, _ = n.Subscribe("service.get", func(msg *nats.Msg) {
		if len(msg.Data) != 0 {
			qs := Service{}
			if err := json.Unmarshal(msg.Data, &qs); err != nil {
				log.Println(err)
			}

			for _, service := range mockServices {
				if qs.GroupID != 0 && service.GroupID == qs.GroupID && service.ID == qs.ID {
					data, _ := json.Marshal(service)
					if err := n.Publish(msg.Reply, data); err != nil {
						log.Println(err)
					}
					return
				} else if qs.GroupID == 0 && service.ID == qs.ID {
					data, _ := json.Marshal(service)
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
}

func findServiceSubscriber() {
	sub, _ := n.Subscribe("service.find", func(msg *nats.Msg) {
		var s []Service
		var qs Service
		if err := json.Unmarshal(msg.Data, &qs); err != nil {
			log.Println(err)
		}

		if qs.Name == "" && qs.ID == "" {
			data, _ := json.Marshal(mockServices)
			if err := n.Publish(msg.Reply, data); err != nil {
				log.Println(err)
			}
			return
		}

		for _, service := range mockServices {
			if service.Name == qs.Name ||
				service.Name == qs.Name && service.Version == qs.Version && qs.GroupID == 0 ||
				service.Name == qs.Name && service.GroupID == qs.GroupID {
				s = append(s, service)
			}
		}

		data, _ := json.Marshal(s)
		if err := n.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func createServiceSubscriber() {
	_, _ = n.Subscribe("service.set", func(msg *nats.Msg) {
		var s Service

		if err := json.Unmarshal(msg.Data, &s); err != nil {
			log.Println(err)
		}
		s.ID = "3"
		data, _ := json.Marshal(s)

		if err := n.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
}

func deleteServiceSubscriber() {
	_, _ = n.Subscribe("service.del", func(msg *nats.Msg) {
		var s Service

		if err := json.Unmarshal(msg.Data, &s); err != nil {
			log.Println(err)
		}

		if err := n.Publish(msg.Reply, []byte{}); err != nil {
			log.Println(err)
		}
	})
}

func notFoundSubscriber(subject string, max int) {
	sub, _ := n.Subscribe(subject, func(msg *nats.Msg) {
		if err := n.Publish(msg.Reply, []byte(`{"_error","Not found"}`)); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(max); err != nil {
		log.Println(err)
	}
}

func foundSubscriber(subject string, resp string, max int) {
	sub, _ := n.Subscribe(subject, func(msg *nats.Msg) {
		if err := n.Publish(msg.Reply, []byte(resp)); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(max); err != nil {
		log.Println(err)
	}
}
