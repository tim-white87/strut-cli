package product

// Template Template used for product file
type Template struct {
	Product Product
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
