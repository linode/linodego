package golinode

import (
	"bytes"
	"fmt"
	"text/template"
)

const (
	stackscriptsName = "stackscripts"
	imagesName       = "images"
	instancesName    = "instances"
	regionsName      = "regions"
	disksName        = "disks"
	configsName      = "configs"
	backupsName      = "backups"
	volumesName      = "volumes"
	kernelsName      = "kernels"
	typesName        = "types"

	stackscriptsEndpoint = "linode/stackscripts"
	imagesEndpoint       = "images"
	instancesEndpoint    = "linode/instances"
	regionsEndpoint      = "regions"
	configsEndpoint      = "linode/instances/{{ .ID }}/configs"
	disksEndpoint        = "linode/instances/{{ .ID }}/disks"
	backupsEndpoint      = "linode/instances/{{ .ID }}/backups"
	volumesEndpoint      = "volumes"
	kernelsEndpoint      = "linode/kernels"
	typesEndpoint        = "linode/types"
)

// Resource represents a linode API resource
type Resource struct {
	name             string
	endpoint         string
	isTemplate       bool
	endpointTemplate *template.Template
}

// NewResource is the factory to create a new Resource struct. If it has a template string the useTemplate bool must be set.
func NewResource(name string, endpoint string, useTemplate bool) *Resource {
	var tmpl *template.Template

	if useTemplate {
		tmpl = template.Must(template.New(name).Parse(endpoint))
	}
	return &Resource{name, endpoint, useTemplate, tmpl}
}

func (r Resource) render(data interface{}) (string, error) {
	if data == nil {
		return "", fmt.Errorf("Cannot template endpoint with <nil> data")
	}
	out := ""
	buf := bytes.NewBufferString(out)
	if err := r.endpointTemplate.Execute(buf, struct{ ID interface{} }{data}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// EndpointWithID will return the rendered endpoint string for the resource with provided id
func (r Resource) EndpointWithID(id int) (string, error) {
	if !r.isTemplate {
		return r.endpoint, nil
	}
	return r.render(id)
}

// Endpoint will return the non-templated endpoint strig for resource
func (r Resource) Endpoint() (string, error) {
	if r.isTemplate {
		return "", fmt.Errorf("Tried to get endpoint for %s without providing data for template", r.name)
	}
	return r.endpoint, nil
}
