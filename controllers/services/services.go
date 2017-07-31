package services

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /services/ with a list of all
// services for current user group
func List(au models.User) (int, []byte) {
	var list []models.Service
	var body []byte
	var user models.User

	users := user.FindAllKeyValue()

	query := make(map[string]interface{}, 0)
	services, err := au.ServicesBy(query)
	if err != nil {
		h.L.Warning(err.Error())
	}
	for _, s := range services {
		exists := false
		for i, e := range list {
			if e.Name == s.Name {
				if e.Version.Before(s.Version) {
					list[i] = s
				}
				exists = true
			}
		}
		if exists == false {
			for id, name := range users {
				if id == s.UserID {
					s.UserName = name
				}
			}
			list = append(list, s)
		}
	}

	body, err = json.Marshal(list)
	if err != nil {
		return 500, []byte("Internal error")
	}

	return http.StatusOK, body
}
