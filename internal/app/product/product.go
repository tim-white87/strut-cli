package product

// Model The product model
type Model interface {
	CreateFile(filetype string)
	ReadFile()
	UpdateFile()
	AddApplication()
}

type model struct {
	*Product
	*FileService
}

// New product model constructor
func New() Model {
	return model{
		&Product{},
		&FileService{},
	}
}

// Product Product anemic model
type Product struct {
	Name         string
	Version      string
	Dependencies []Dependency
	Applications []Application
}

// Dependency Product tool/language/installation dependency
type Dependency struct{}

// Application Application defined in the product
type Application struct{}

// FileService Manages the product file
type FileService struct{}

// CreateFile Creates the product file in JSON or YAML
func (fs *FileService) CreateFile(filetype string) {}

// ReadFile Loads the product file from the CWD
func (fs *FileService) ReadFile() {}

// UpdateFile Updates the product file
func (fs *FileService) UpdateFile() {}

// AddApplication Adds an application to product file
func (fs *FileService) AddApplication() {}
