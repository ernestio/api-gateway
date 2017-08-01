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
	var s models.Service
	var err error

	if _, ok := query["name"]; !ok {
		return 500, []byte("Internal error")
	}
	name := query["name"].(string)

	if !au.IsReader(s.GetType(), name) {
		return 403, []byte("You're not allowed to access this resource")
	}

	if s, err = s.FindLastByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	if s.ID == "" {
		return 404, []byte("Specified environment name does not exist")
	}

	if err := o.Render(s); err != nil {
		h.L.Warning(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}
	if body, err = o.ToJSON(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, body
}
