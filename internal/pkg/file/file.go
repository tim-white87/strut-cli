package file

// Manager CRUD operations on a file
type Manager interface {
	CreateFile()
	ReadFile()
	UpdateFile()
	DeleteFile()
}
