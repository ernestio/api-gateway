package services

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : Not implemented
func Update(au models.User, name string, body []byte) (int, []byte) {
	var raw []byte
	var err error
	var input models.Service

	if err := json.Unmarshal(body, &input); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Get existing service
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

	s.Sync = input.Sync
	s.SyncType = input.SyncType
	s.SyncInterval = input.SyncInterval
	if s.Sync == true {
		if s.SyncType != "hard" {
			s.SyncType = "soft"
		}
		if s.SyncInterval == 0 {
			s.SyncInterval = 5
		}
	}

	if err := s.Save(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"id":"` + s.ID + `"}`)
}
