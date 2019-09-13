package provider

// Types provider types
var Types = &providerRegistry{
	AWS: &Type{"AWS"},
}

// TypeList list of types
var TypeList = []*Type{Types.AWS}

type providerRegistry struct {
	AWS *Type
}

// Type Provider type
type Type struct {
	Name string
}
