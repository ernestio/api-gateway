package datacenters

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Get : responds to GET /datacenter/:id:/ with the specified
// datacenter details
func Get(au models.User, datacenter string) (int, []byte) {
	var s models.Service
	var envs []models.Service
	var r models.Role
	var roles []models.Role
	var d models.Datacenter
	var body []byte
	var err error

	if err := d.FindByName(datacenter, &d); err != nil {
		return 404, []byte("Project not found")
	}

	query := make(map[string]interface{})
	query["datacenter_id"] = d.ID
	if err := s.Find(query, &envs); err == nil {
		for _, v := range envs {
			d.Environments = append(d.Environments, v.Name)
		}
	}

	if err := r.FindAllByResource(d.GetID(), d.GetType(), &roles); err == nil {
		for _, v := range roles {
			d.Roles = append(d.Roles, v.UserID+" ("+v.Role+")")
		}
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
