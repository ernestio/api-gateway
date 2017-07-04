package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
	"github.com/nu7hatch/gouuid"
)

func getServicesOutput(filter map[string]interface{}) (list []views.ServiceRender, err error) {
	var s models.Service
	var services []models.Service
	var o views.ServiceRender

	if err := s.Find(filter, &services); err != nil {
		return list, err
	}

	return o.RenderCollection(services)
}

func getServiceRaw(name string, group int) (service []byte, err error) {
	var s models.Service
	var services []models.Service

	if err = s.FindByNameAndGroupID(name, group, &services); err != nil {
		return nil, errors.New(`"Internal error"`)
	}

	if len(services) == 0 {
		return nil, errors.New(`"Service not found"`)
	}

	body, err := json.Marshal(services[0])
	if err != nil {
		return nil, errors.New("Internal error")
	}
	return body, nil
}

func generateStreamID(salt string) string {
	compose := []byte(salt)
	hasher := md5.New()
	if _, err := hasher.Write(compose); err != nil {
		h.L.Warning(err.Error())
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func getDatacenter(name string, group int) (datacenter []byte, err error) {
	var d models.Datacenter
	var datacenters []models.Datacenter

	if err := d.FindByNameAndGroupID(name, group, &datacenters); err != nil {
		return datacenter, err
	}

	if len(datacenters) == 0 {
		return datacenter, errors.New(`"Specified datacenter does not exist"`)
	}

	datacenter, err = json.Marshal(datacenters[0])
	if err != nil {
		return datacenter, errors.New("Internal error trying to get the datacenter")
	}

	return datacenter, nil
}

func getGroup(id int) (group []byte, err error) {
	var g models.Group

	if err = g.FindByID(id); err != nil {
		return group, errors.New(`"Specified group does not exist"`)
	}

	if group, err = json.Marshal(g); err != nil {
		return group, errors.New(`"Internal error"`)
	}
	h.L.Info(group)

	return group, nil
}

// Generates a service id composed by a random uuid, and
// a valid generated stream id
func generateServiceID(salt string) string {
	sufix := generateStreamID(salt)
	prefix, _ := uuid.NewV4()

	return prefix.String() + "-" + string(sufix[:])
}

func getService(name string, group int) (service *models.Service, err error) {
	var s models.Service
	var services []models.Service

	if err = s.FindByNameAndGroupID(name, group, &services); err != nil {
		return service, h.ErrGatewayTimeout
	}

	if len(services) == 0 {
		return nil, nil
	}

	return &services[0], nil
}
