package services

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// ForceDeletion : Deletes a service by name forcing it
func ForceDeletion(au models.User, name string) (int, []byte) {
	var service models.Service

	if !au.IsOwner(service.GetType(), name) {
		return 403, []byte("You're not allowed to access this resource")
	}

	if err := service.DeleteByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"id":"` + service.ID + `"}`)
}
