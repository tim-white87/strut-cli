package file

// Types file types
var Types = &TypeRegistry{
	JSON: "json",
	YAML: "yaml",
}

// TypeRegistry Struct for various file types
type TypeRegistry struct {
	JSON string
	YAML string
}

// Manager CRUD operations on a file
type Manager interface {
	CreateFile()
	ReadFile()
	UpdateFile()
	DeleteFile()
}
