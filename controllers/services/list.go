package services

import (
	"encoding/json"
	"net/http"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /services/ with a list of all
// services for current user group
func List(au models.User) (int, []byte) {
	var list []models.Service
	var body []byte
	var user models.User

	query := make(map[string]interface{}, 0)
	services, err := au.ServicesBy(query)
	if err != nil {
		h.L.Warning(err.Error())
		return 404, []byte("Environment not found")
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
			for id, name := range user.FindAllKeyValue() {
				if id == s.UserID {
					s.UserName = name
				}
			}
			var d models.Datacenter
			var ds []models.Datacenter
			if err := d.FindAll(au, &ds); err == nil {
				for _, d = range ds {
					s.Project = d.Name
				}
			}
			list = append(list, s)
		}
	}

	for i := range list {
		nameParts := strings.Split(list[i].Name, models.EnvNameSeparator)
		list[i].Name = nameParts[1]
	}

	body, err = json.Marshal(list)
	if err != nil {
		return 500, []byte("Internal error")
	}

	return http.StatusOK, body
}
