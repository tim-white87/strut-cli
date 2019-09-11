package file

// Types file types
var Types = &TypeRegistry{
	JSON: &Type{"json"},
	YAML: &Type{"yaml"},
}

// TypeList list of types
var TypeList = []*Type{Types.JSON, Types.YAML}

// TypeRegistry Struct for various file types
type TypeRegistry struct {
	JSON *Type
	YAML *Type
}

// Type file type
type Type struct {
	Extension string
}
