package envs

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
	var input models.Env

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateEnv, input.GetType(), name); st != 200 {
		return st, res
	}

	if err := json.Unmarshal(body, &input); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Get existing environment
	if raw, err = getEnvRaw(au, name); err != nil {
		return 404, []byte(err.Error())
	}

	s := models.Env{}
	if err := json.Unmarshal(raw, &s); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	if s.Status == "in_progress" {
		return 400, []byte(`"Environment is already applying some changes, please wait until they are done"`)
	}

	s.Options["sync"] = input.Options["sync"]
	s.Options["sync_type"] = input.Options["sync_type"]
	s.Options["sync_interval"] = input.Options["sync_interval"]
	if s.Options["sync"] == true {
		if s.Options["sync_type"] != "hard" {
			s.Options["sync_type"] = "soft"
		}
		if s.Options["sync_interval"] == 0 {
			s.Options["sync_interval"] = 5
		}
	}

	if err := s.Save(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"id":"` + s.ID + `"}`)
}
