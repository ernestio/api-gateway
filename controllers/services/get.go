package services

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Get : responds to GET /services/:service with the
// details of an existing service
func Get(au models.User, query map[string]interface{}) (int, []byte) {
	var o views.ServiceRender
	var body []byte

	services, err := au.ServicesBy(query)
	if err != nil {
		return 500, []byte(err.Error())
	}

	if len(services) > 0 {
		if err := o.Render(services[0]); err != nil {
			h.L.Warning(err.Error())
			return http.StatusBadRequest, []byte(err.Error())
		}
		if body, err = o.ToJSON(); err != nil {
			return 500, []byte(err.Error())
		}
		return http.StatusOK, body
	}

	return http.StatusNotFound, []byte("")
}
