package provider

// Types provider types
var Types = &providerRegistry{
	AWS: &Type{Name: "AWS"},
}

// TypeList list of types
var TypeList = []*Type{Types.AWS}

type providerRegistry struct {
	AWS *Type
}

// Type Provider type
type Type struct {
	Name  string
	Model Model `json:"-"`
}

// Provider Hosted application provider
type Provider struct {
	Type      *Type       `json:"type"`
	Resources *[]Resource `json:"resources"`
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
	Load(*Provider)
	Provision()
	Destroy()
	CheckStatus()
}
