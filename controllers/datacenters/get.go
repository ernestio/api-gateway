package datacenters

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ernestio/api-gateway/models"
)

// Get : responds to GET /datacenter/:id:/ with the specified
// datacenter details
func Get(datacenter string) (int, []byte) {
	var d models.Datacenter
	var body []byte
	var err error

	id, err := strconv.Atoi(datacenter)
	if err != nil {
		return 404, []byte("Datacenter not found")
	}
	if err := d.FindByID(id); err != nil {
		return 404, []byte("Datacenter not found")
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
