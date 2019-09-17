package provider

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

// Provider Hosted application provider
type Provider struct {
	Name      string      `json:"name"`
	Resources []*Resource `json:"resources"`
	*ResourceCommands
}

// Resource Provider resource (i.e. cloudformation)
type Resource struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Body     string `json:"body,omitempty"`
	Priority int    `json:"priority,omitempty"`
	*ResourceCommands
}

// ResourceCommands Custom resource commands
type ResourceCommands struct {
	PreProvision  []string `json:"preProvision"`
	PostProvision []string `json:"postProvision"`
}

// Model provider model interface
type Model interface {
	Provision([]*Resource)
	Destroy()
	CheckStatus()
}

// ProvisionResources initiates resource provisioning of provider map
func ProvisionResources(provisionMap map[string]map[int][]*Resource) {
	for provider, resourcesMap := range provisionMap {
		model := ModelsMap[provider]
		// TODO batch each priority group with concurrency
		// model.Provision(resources)

		// TODO run 0 items last
		model.Provision(resourcesMap[0])
	}
}
