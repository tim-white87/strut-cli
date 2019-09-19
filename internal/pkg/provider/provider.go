package provider

import (
	"sync"

	"github.com/fatih/color"
)

// ModelsMap Maps provider name to model
var ModelsMap = map[string]Model{
	Types.AWS: NewAwsModel(),
}

// Types provider types
var Types = &providerRegistry{
	AWS: "AWS",
}

type providerRegistry struct {
	AWS string
}

// Resource Provider resource (i.e. cloudformation)
type Resource struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Body     string    `json:"body,omitempty"`
	Priority int       `json:"priority,omitempty"`
	Provider *Provider `json:"provider"`
	*ResourceCommands
}

// Provider Hosted application provider
type Provider struct {
	Name string `json:"name"`
	*ResourceCommands
}

// ResourceCommands Custom resource commands
type ResourceCommands struct {
	PreProvision  []string `json:"preProvision"`
	PostProvision []string `json:"postProvision"`
}

// Model provider model interface
type Model interface {
	Provision(*Resource)
	Destroy(*Resource)
	CheckStatus()
}

// Provision initiates resource provisioning of provider map
func Provision(provisionMap map[int][]*Resource) {
	batchWg := &sync.WaitGroup{}
	batchWg.Add(len(provisionMap))
	defer batchWg.Wait()

	for priority, batch := range provisionMap {
		color.Green("Batch #: %b", priority)
		go provisionBatch(batch, batchWg)
	}
}

func provisionBatch(batch []*Resource, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, resource := range batch {
		provisionResource(resource)
	}
}

func provisionResource(r *Resource) {
	color.Green("Provisioning >>> Resource: %s on Provider: %s >>> ", r.Name, r.Provider.Name)
	model := ModelsMap[r.Provider.Name]
	model.Provision(r)
}
