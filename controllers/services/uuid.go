package services

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// UUID : Creates an unique id
func UUID(au models.User, body []byte) (int, []byte) {
	var s struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &s); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}
	id := generateStreamID(s.ID)

	return http.StatusOK, []byte(`{"uuid":"` + id + `"}`)
}
