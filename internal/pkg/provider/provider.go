package provider

import (
	"sort"
	"sync"

	"github.com/fatih/color"
)

// ModelsMap Maps provider name to model
var ModelsMap = map[string]func(*Resource) Model{
	Types.AWS: NewAwsModel,
}

// Types provider types
var Types = &providerRegistry{
	AWS: "AWS",
}

type providerRegistry struct {
	AWS string
}

// Status status
var Status = &statusRegistry{
	"Complete",
	"InProgress",
	"Failed",
	"NotFound",
}

type statusRegistry struct {
	Complete   string
	InProgress string
	Failed     string
	NotFound   string
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
	Provision()
	Destroy()
	CheckStatus() string
}

func sortProvisionMap(provisionMap map[int][]*Resource) []int {
	keys := make([]int, 0)
	for k := range provisionMap {
		if k != 0 {
			keys = append(keys, k)
		}
	}
	sort.Ints(keys)
	keys = append(keys, 0)
	return keys
}

func reverseInts(input []int) []int {
	if len(input) == 0 {
		return input
	}
	return append(reverseInts(input[1:]), input[0])
}

// Provision initiates resource provisioning of provider map
func Provision(provisionMap map[int][]*Resource) {
	keys := sortProvisionMap(provisionMap)
	for _, priority := range keys {
		batch := provisionMap[priority]
		provisionBatch(batch, priority)
	}
}

func provisionBatch(batch []*Resource, priority int) {
	if priority == 0 {
		color.HiBlack("Batch >>> Final")
	} else {
		color.HiBlack("Batch >>> Priority: #%b", priority)
	}
	resourceWg := &sync.WaitGroup{}
	resourceWg.Add(len(batch))
	defer resourceWg.Wait()
	for _, resource := range batch {
		go provisionResource(resource, resourceWg)
	}
}

func provisionResource(r *Resource, wg *sync.WaitGroup) {
	defer wg.Done()
	model := ModelsMap[r.Provider.Name](r)
	model.Provision()
}

// Destroy blows up the infrastructure in the cloud
func Destroy(provisionMap map[int][]*Resource) {
	keys := reverseInts(sortProvisionMap(provisionMap))
	for _, priority := range keys {
		batch := provisionMap[priority]
		destroyBatch(batch, priority)
	}
}

func destroyBatch(batch []*Resource, priority int) {
	if priority == 0 {
		color.HiBlack("Batch >>> Final")
	} else {
		color.HiBlack("Batch >>> Priority: #%b", priority)
	}
	resourceWg := &sync.WaitGroup{}
	resourceWg.Add(len(batch))
	defer resourceWg.Wait()
	for _, resource := range batch {
		go destroyResource(resource, resourceWg)
	}
}

func destroyResource(r *Resource, wg *sync.WaitGroup) {
	defer wg.Done()
	model := ModelsMap[r.Provider.Name](r)
	model.Destroy()
}
