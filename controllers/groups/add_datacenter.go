package groups

import (
	"encoding/json"
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// AddDatacenter : Adds a datacenter to a group
func AddDatacenter(au models.User, g string, body []byte) (int, []byte) {
	var group models.Group
	var datacenter models.Datacenter
	var payload map[string]string

	groupID, err := strconv.Atoi(g)
	if err != nil {
		h.L.Error("Invalid input adding a datacenter to a group")
		return http.StatusBadRequest, []byte("Invalid input")
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		h.L.Error("Error unmarshalling input to a datacenter struct")
		return http.StatusBadRequest, []byte("Invalid input")
	}

	datacenterID, err := strconv.Atoi(payload["datacenterid"])
	if err != nil {
		h.L.Error("Provided datacenter identifier is not an integer")
		return http.StatusBadRequest, []byte("Invalid input")
	}

	if err := group.FindByID(groupID); err != nil {
		h.L.Error("Couldn't found specified group id")
		return 500, []byte("Internal server error")
	}

	if err := datacenter.FindByID(datacenterID); err != nil {
		h.L.Error("Couldn't found specified datacenter id")
		return 500, []byte("Internal server error")
	}

	datacenter.GroupID = groupID
	if err = datacenter.Save(); err != nil {
		h.L.Error("Provided datacenter does not belong to given group")
		return http.StatusBadRequest, []byte("Provided datacenter does not belong to given group")
	}

	return http.StatusOK, []byte("Datacenter successfully added to group " + group.Name)
}
