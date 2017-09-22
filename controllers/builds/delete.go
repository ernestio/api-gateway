package builds

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : Deletes an environment by name, generating a delete build
func Delete(au models.User, name string) (int, []byte) {
	var e models.Env
	var m models.Mapping

	err := e.FindByName(name)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the environment"`)
	}

	err = m.Delete(name)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the environment"`)
	}

	b := models.Build{
		EnvironmentID: e.ID,
		UserID:        au.ID,
		Username:      au.Username,
		Type:          "destroy",
		Mapping:       m,
	}

	err = b.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't create the deletion build"`)
	}

	if err := b.RequestDeletion(&m); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't call service.delete"`)
	}

	return http.StatusOK, []byte(`{"id":"` + b.ID + `"}`)
}
