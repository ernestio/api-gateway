package builds

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Get : responds to GET /services/:service with the
// details of an existing service
func Get(au models.User, id string) (int, []byte) {
	var o views.BuildRender
	var err error
	var body []byte
	var e models.Env
	var b models.Build
	var p models.Project

	if err = b.FindByID(id); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	if b.ID == "" {
		return 404, []byte("Specified environment name does not exist")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, b.GetType(), id); st != 200 {
		return st, res
	}

	if err := e.FindByID(int(b.EnvironmentID)); err != nil {
		return 404, []byte("Environment not found")
	}

	if err := p.FindByID(int(e.ProjectID)); err != nil {
		return 404, []byte("Environment not found")
	}

	o.Name = e.Name
	o.Project = p.Name
	o.Provider = e.Type

	if err := o.Render(b); err != nil {
		h.L.Warning(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}
	if body, err = o.ToJSON(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, body
}
