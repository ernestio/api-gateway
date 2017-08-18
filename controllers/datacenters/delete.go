package datacenters

import (
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /datacenters/:id: by deleting an
// existing datacenter
func Delete(au models.User, datacenter string) (int, []byte) {
	var d models.Datacenter

	id, err := strconv.Atoi(datacenter)
	if err = d.FindByID(id); err != nil {
		return 404, []byte("Datacenter not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.DeleteProject, d.GetType(), d.Name); st != 200 {
		return st, res
	}

	ss, err := d.Services()
	if err != nil {
		return 500, []byte(err.Error())
	}

	if len(ss) > 0 {
		return 400, []byte("Existing environments are referring to this project.")
	}

	if err := d.Delete(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte("Datacenter successfully deleted")
}
