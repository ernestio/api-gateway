package envs

import (
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Reset : Respons to POST /services/:service/reset/ and updates the
// service status to errored from in_progress
func Reset(au models.User, name string) (int, []byte) {
	var s models.Env
	var envs []models.Env

	if st, res := h.IsAuthorizedToResource(&au, h.ResetBuild, s.GetType(), name); st != 200 {
		return st, res
	}

	filter := make(map[string]interface{})
	filter["name"] = name
	if err := s.Find(filter, &envs); err != nil {
		h.L.Warning(err.Error())
		return 500, []byte("Internal Error")
	}

	if len(envs) == 0 {
		return 404, []byte("Environment not found with this name")
	}

	s = envs[0]

	if s.Status != "in_progress" {
		return 200, []byte("Reset only applies to an 'in progress' environment, however environment '" + name + "' is on status '" + s.Status)
	}

	if err := s.Reset(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	return 200, []byte("success")
}
