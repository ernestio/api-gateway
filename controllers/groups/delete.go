package groups

import (
	"net/http"
	"strconv"

	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /groups/:id: by deleting an
// existing group
func Delete(au models.User, group string) (int, []byte) {
	var g models.Group
	var users []models.User
	var datacenters []models.Datacenter

	if au.Admin != true {
		return 403, []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action")
	}

	id, err := strconv.Atoi(group)
	if err = g.FindByID(id); err != nil {
		return 404, []byte("Specified group does not exists")
	}

	// Check if there are any users on the group
	if users, err = g.Users(); err != nil {
		return 500, []byte("Internal server error")
	}

	if len(users) > 0 {
		return 400, []byte("This group has users assigned to it, please remove the users before performing this action")
	}

	// Check if there are any datacenters on the group
	if datacenters, err = g.Datacenters(); err != nil {
		return 500, []byte("Internal server error")
	}

	if len(datacenters) > 0 {
		return 400, []byte("This group has datacenters assigned to it, please remove the datacenters before performing this action")
	}

	if err := g.Delete(); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, []byte("Group successfully deleted")
}
