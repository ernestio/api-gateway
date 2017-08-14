package services

import (
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : Deletes a service by name
func Delete(au models.User, name string) (int, []byte) {
	var err error
	var def models.Definition
	var s models.Service

	if !au.IsOwner(s.GetType(), name) {
		return 403, []byte("You're not allowed to access this resource")
	}

	if s, err = s.FindLastByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	if s.ID == "" {
		return 404, []byte("Specified environment name does not exist")
	}

	if s.Status == "in_progress" {
		return 400, []byte(`"Service is already applying some changes, please wait until they are done"`)
	}

	dID := strconv.Itoa(s.DatacenterID)
	body, err := def.MapDeletion(s.ID, s.Type, dID)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the service"`)
	}
	if err := s.RequestDeletion(body); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't call service.delete"`)
	}

	return http.StatusOK, []byte(`{"id":"` + s.ID + `"}`)
}
