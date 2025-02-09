package outlet

import (
	"errors"
)

// Outlet interface for all smart outlets
type Outlet interface {
	getID() int
	getBrand() string
	action(action string) error
}

// NewOutlet factory function to create an outlet based on the brand
func newOutlet(brand string, id int) (Outlet, error) {
	switch brand {
	case "kasa":
		return &kasaOutlet{id: id}, nil
	default:
		return nil, errors.New("unsupported outlet brand")
	}
}
