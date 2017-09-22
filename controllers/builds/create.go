package builds

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
	"github.com/ernestio/mapping/definition"
)

// Create : Creates an environment build
func Create(au models.User, definition *definition.Definition, raw []byte, dry string) (int, []byte) {
	var e models.Env
	var m models.Mapping

	err := e.FindByName(definition.FullName())
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the environment"`)
	}

	err = m.Apply(definition)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the environment"`)
	}

	if dry == "true" {
		res, err := views.RenderDefinition(m)
		if err != nil {
			h.L.Error(err.Error())
			return 400, []byte("Internal error")
		}
		return http.StatusOK, res
	}

	b := models.Build{
		ID:            m["id"].(string),
		EnvironmentID: e.ID,
		UserID:        au.ID,
		Username:      au.Username,
		Type:          "apply",
		Mapping:       m,
	}

	err = b.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't create the build"`)
	}

	if err := b.RequestCreation(&m); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't call service.create"`)
	}

	return http.StatusOK, []byte(`{"id":"` + b.ID + `"}`)
}
