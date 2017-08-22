package envs

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// ForceDeletion : Deletes a service by name forcing it
func ForceDeletion(au models.User, name string) (int, []byte) {
	var s models.Env

	if st, res := h.IsAuthorizedToResource(&au, h.DeleteEnvForce, s.GetType(), name); st != 200 {
		return st, res
	}

	if err := s.DeleteByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"id":"` + s.ID + `"}`)
}
