package services

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Search : Finds all services
func Search(au models.User, query map[string]interface{}) (int, []byte) {
	var o views.ServiceRender

	services, err := au.ServicesBy(query)
	if err != nil {
		return 500, []byte(err.Error())
	}

	list, err := o.RenderCollection(services)
	if err != nil {
		return 500, []byte(err.Error())
	}

	b, err := json.Marshal(list)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	return http.StatusOK, b
}
