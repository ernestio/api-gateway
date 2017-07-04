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
	var err error
	var s models.Service
	var services []models.Service
	var o views.ServiceRender
	var body []byte

	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if err = s.Find(query, &services); err != nil {
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
