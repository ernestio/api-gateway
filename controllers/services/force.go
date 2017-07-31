package services

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// ForceDeletion : Deletes a service by name forcing it
func ForceDeletion(au models.User, name string) (int, []byte) {
	var raw []byte
	var err error
	var service models.Service

	if raw, err = getServiceRaw(au, name); err != nil {
		return 404, []byte(err.Error())
	}

	s := models.Service{}
	if err := json.Unmarshal(raw, &s); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(err.Error())
	}

	if err := service.DeleteByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"id":"` + s.ID + `"}`)
}
