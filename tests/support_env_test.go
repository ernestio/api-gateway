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
	mockServices = []models.Env{
		models.Env{
			ID:   1,
			Name: "fake-test",
		},
		models.Env{
			ID:   3,
			Name: "fake-test",
		},
		models.Env{
			ID:   2,
			Name: "fake-test2",
		},
	}
)

func findServiceSolo() {
	sub, _ := models.N.Subscribe("environment.find", func(msg *nats.Msg) {
		var s []models.Env
		var qs models.Env
		if err := json.Unmarshal(msg.Data, &qs); err != nil {
			log.Println(err)
		}

		if qs.Name == "" && qs.ID == 0 {
			data, _ := json.Marshal(mockServices)
			if err := models.N.Publish(msg.Reply, data); err != nil {
				log.Println(err)
			}
			return
		}

		for _, service := range mockServices {
			if service.Name == qs.Name {
				s = append(s, service)
			}
		}

		data, _ := json.Marshal(s)
		if err := models.N.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func findService() {
	sub2, _ := models.N.Subscribe("authorization.find", func(msg *nats.Msg) {
		res := `[{"resource_id":"` + mockServices[0].Name + `"},{"resource_id":"` + mockServices[1].Name + `"}]`
		if err := models.N.Publish(msg.Reply, []byte(res)); err != nil {
			log.Println(err)
		}
	})
	if err := sub2.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func createServiceSubscriber() {
	_, _ = models.N.Subscribe("environment.set", func(msg *nats.Msg) {
		var s models.Env

		if err := json.Unmarshal(msg.Data, &s); err != nil {
			log.Println(err)
		}
		s.ID = 3
		data, _ := json.Marshal(s)

		if err := models.N.Publish(msg.Reply, data); err != nil {
			log.Println(err)
		}
	})

	sub2, _ := models.N.Subscribe("authorization.set", func(msg *nats.Msg) {
		res := `{}`
		if err := models.N.Publish(msg.Reply, []byte(res)); err != nil {
			log.Println(err)
		}
	})
	if err := sub2.AutoUnsubscribe(1); err != nil {
		log.Println(err)
	}
}

func serviceResetSubscriber() {
	_, _ = models.N.Subscribe("build.set.status", func(msg *nats.Msg) {
		if err := models.N.Publish(msg.Reply, []byte(`{"status":"success"}`)); err != nil {
			log.Println(err)
		}
	})
}

func notFoundSubscriber(subject string, max int) {
	sub, _ := models.N.Subscribe(subject, func(msg *nats.Msg) {
		if err := models.N.Publish(msg.Reply, []byte(`{"_error","Not found"}`)); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(max); err != nil {
		log.Println(err)
	}
}

func foundSubscriber(subject string, resp string, max int) {
	sub, _ := models.N.Subscribe(subject, func(msg *nats.Msg) {
		if err := models.N.Publish(msg.Reply, []byte(resp)); err != nil {
			log.Println(err)
		}
	})
	if err := sub.AutoUnsubscribe(max); err != nil {
		log.Println(err)
	}
}
