package services

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Sync : Respons to POST /services/:service/sync/ and synchronizes a service with
// its provider representation
func Sync(au models.User, name string) (int, []byte) {
	var raw []byte
	var err error

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

	if err = s.RequestSync(); err != nil {
		return 500, []byte("An error ocurred while ernest was trying to sync your service")
	}

	// TODO : This probably needs to use the monit tool instead of this.

	return http.StatusOK, []byte("....")
}
