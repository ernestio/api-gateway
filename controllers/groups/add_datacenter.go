package groups

import (
	"encoding/json"
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// AddDatacenter : Adds a datacenter to a group
func AddDatacenter(au models.User, g, d string, body []byte) (int, []byte) {
	var group models.Group
	var datacenter models.Datacenter
	var payload map[string]string

	if au.Admin != true {
		return http.StatusForbidden, []byte("You don't have permissions to perform this action, please login with an admin account")
	}

	groupID, err := strconv.Atoi(g)
	if err != nil {
		return http.StatusBadRequest, []byte("Invalid input")
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		return http.StatusBadRequest, []byte("Invalid input")
	}

	datacenterID, err := strconv.Atoi(d)
	if err != nil {
		return http.StatusBadRequest, []byte("Invalid input")
	}

	if err := group.FindByID(groupID); err != nil {
		return 500, []byte("Internal server error")
	}

	if err := datacenter.FindByID(datacenterID); err != nil {
		return 500, []byte("Internal server error")
	}

	datacenter.GroupID = groupID
	if err = datacenter.Save(); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte("Provided datacenter does not belong to given group")
	}

	return http.StatusOK, []byte("Datacenter successfully added to group " + group.Name)
}
