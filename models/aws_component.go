package models

import (
	"errors"
)

// AWSComponent : AWS Component generic representation
type AWSComponent struct {
	Datacenter *Project          `json:"datacenter"`
	Tags       map[string]string `json:"tags"`
	Name       string            `json:"name"`
}

// FindBy : Will find by the data specified on the component
func (g *AWSComponent) FindBy() (list []interface{}, err error) {
	d := g.Datacenter
	component := g.Name
	tags := g.Tags

	query := make(map[string]interface{})
	query["expects_response"] = true
	query["aws_access_key_id"] = d.AccessKeyID
	query["aws_secret_access_key"] = d.SecretAccessKey
	query["datacenter_region"] = d.Region

	if len(tags) > 0 {
		query["tags"] = tags
	}
	components := make(map[string]interface{})
	if err := NewBaseModel(component).CallStoreBy("find.aws", query, &components); err != nil {
		return list, errors.New("Internal error occurred")
	}

	if components["components"] == nil {
		return list, nil
	}

	list = components["components"].([]interface{})
	for i := range list {
		component := list[i].(map[string]interface{})
		delete(component, "_batch_id")
		delete(component, "_type")
		delete(component, "aws_access_key_id")
		delete(component, "aws_secret_access_key")
		delete(component, "_uuid")
	}
	return list, nil
}
