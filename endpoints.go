package golinode

import (
	"bytes"
	"fmt"
	"text/template"
)

const (
	stackscriptsEndpoint  = "linode/stackscripts"
	distributionsEndpoint = "linode/distributions"
	instancesEndpoint     = "linode/instances"
	regionsEndpoint       = "regions"
	backupsEndpoint       = "linode/instances/{{ .ID }}/backups"
)

var endpoints = map[string]*Resource{
	"stackscripts":  NewResource("stackscripts", stackscriptsEndpoint, false),
	"distributions": NewResource("distributions", distributionsEndpoint, false),
	"instances":     NewResource("instances", instancesEndpoint, false),
	"regions":       NewResource("regions", regionsEndpoint, false),
	"backups":       NewResource("backups", backupsEndpoint, true),
}

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

// Endpoint will return the endpoint string for the resource, optionally send data if the endpoint is templated
func (r Resource) Endpoint(data interface{}) (string, error) {
	if !r.isTemplate {
		return r.endpoint, nil
	}
	return r.render(data)
}

func (c Client) getResource(resource string) (*Resource, error) {
	selectedResource, ok := endpoints[resource]
	if !ok {
		return nil, fmt.Errorf("Could not find resource %s", resource)
	}
	return selectedResource, nil
}
