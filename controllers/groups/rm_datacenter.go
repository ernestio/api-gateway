package groups

import (
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// RmDatacenter : Deletes a datacenter from a group
func RmDatacenter(au models.User, g, d string) (int, []byte) {
	var group models.Group
	var datacenter models.Datacenter

	if au.Admin != true {
		return http.StatusForbidden, []byte("You don't have permissions to perform this action, please login with an admin account")
	}

	groupid, err := strconv.Atoi(g)
	if err != nil {
		return http.StatusBadRequest, []byte("Invalid input")
	}
	if err = group.FindByID(groupid); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	datacenterid, err := strconv.Atoi(d)
	if err = datacenter.FindByID(datacenterid); err != nil {
		return http.StatusGatewayTimeout, []byte("Internal server error")
	}

	datacenter.GroupID = 0
	if err = datacenter.Save(); err != nil {
		h.L.Error(err.Error())
		return http.StatusGatewayTimeout, []byte("Internal server error")
	}

	return http.StatusOK, []byte("Datacenter successfully removed from group " + group.Name)
}
