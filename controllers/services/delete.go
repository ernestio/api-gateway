package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : Deletes a service by name
func Delete(au models.User, name string) (int, []byte) {
	var raw []byte
	var err error
	var def models.Definition

	if raw, err = getServiceRaw(name, au.GroupID); err != nil {
		return 404, []byte(err.Error())
	}

	s := models.Service{}
	if err := json.Unmarshal(raw, &s); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	if s.Status == "in_progress" {
		return 400, []byte(`"Service is already applying some changes, please wait until they are done"`)
	}

	dID := strconv.Itoa(s.DatacenterID)
	body, err := def.MapDeletion(s.ID, s.Type, dID)
	if err != nil {
		return 500, []byte(`"Couldn't map the service"`)
	}
	if err := s.RequestDeletion(body); err != nil {
		return 500, []byte(`"Couldn't call service.delete"`)
	}

	parts := strings.Split(s.ID, "-")
	stream := parts[len(parts)-1]

	return http.StatusOK, []byte(`{"id":"` + s.ID + `","stream_id":"` + stream + `"}`)
}
