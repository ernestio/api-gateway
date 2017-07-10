package services

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Builds : gets the list of builds for the specified service
func Builds(au models.User, query map[string]interface{}) (int, []byte) {
	var o views.ServiceRender

	builds, err := au.ServicesBy(query)
	if err != nil {
		return 500, []byte(err.Error())
	}

	list, err := o.RenderCollection(builds)
	if err != nil {
		return 500, []byte(err.Error())
	}

	users := au.FindAllKeyValue()
	for i := range list {
		for id, name := range users {
			if id == list[i].UserID {
				list[i].UserName = name
			}
		}
	}

	body, err := json.Marshal(list)
	if err != nil {
		h.L.Warning(err.Error())
		return 500, []byte("Internal error")
	}

	return http.StatusOK, body
}
