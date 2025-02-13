package outlet

import (
	"errors"
)

type Outlet interface {
	getID() string
	getBrand() string
	action(action string) error
}

func newOutlet(brand string, id string) (Outlet, error) {
	switch brand {
	case "kasa":
		return &kasaOutlet{id: id}, nil
	default:
		return nil, errors.New("unsupported outlet brand")
	}
}
