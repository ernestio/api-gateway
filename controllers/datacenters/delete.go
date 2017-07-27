package datacenters

import (
	"net/http"
	"strconv"

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

	if ok := au.Owns(&d); !ok {
		return http.StatusForbidden, []byte("You don't have permissions to acccess this resource")
	}

	ss, err := d.Services()
	if err != nil {
		return 500, []byte(err.Error())
	}

	if len(ss) > 0 {
		return 400, []byte("Existing services are referring to this datacenter.")
	}

	if err := d.Delete(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte("Datacenter successfully deleted")
}
