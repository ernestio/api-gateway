package datacenters

import (
	"encoding/json"
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : responds to PUT /datacenters/:id: by updating
// an existing datacenter
func Update(au models.User, datacenter string, body []byte) (int, []byte) {
	var d models.Datacenter
	var existing models.Datacenter
	var err error

	if d.Map(body) != nil {
		return 400, []byte("Invalid input")
	}

	id, err := strconv.Atoi(datacenter)
	if err = existing.FindByID(id); err != nil {
		return 404, []byte("Datacenter not found")
	}

	if ok := au.Owns(&d); !ok {
		return http.StatusForbidden, []byte("You don't have permissions to acccess this resource")
	}

	existing.Username = d.Username
	existing.Password = d.Password
	existing.AccessKeyID = d.AccessKeyID
	existing.SecretAccessKey = d.SecretAccessKey
	existing.SubscriptionID = d.SubscriptionID
	existing.ClientID = d.ClientID
	existing.ClientSecret = d.ClientSecret
	existing.TenantID = d.TenantID
	existing.Environment = d.Environment

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
