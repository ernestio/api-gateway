package services

import (
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Reset : Respons to POST /services/:service/reset/ and updates the
// service status to errored from in_progress
func Reset(au models.User, name string) (int, []byte) {
	var s models.Service
	var services []models.Service

	filter := make(map[string]interface{})
	filter["name"] = name
	if err := s.Find(filter, &services); err != nil {
		h.L.Warning(err.Error())
		return 500, []byte("Internal Error")
	}

	if len(services) == 0 {
		return 404, []byte("Service not found with this name")
	}

	s = services[0]

	if s.Status != "in_progress" {
		return 200, []byte("Reset only applies to 'in progress' serices, however service '" + name + "' is on status '" + s.Status)
	}

	if err := s.Reset(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	return 200, []byte("success")
}
