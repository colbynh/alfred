package devices

import (
	"errors"
	"fmt"
)

type Device interface {
	GetClass() string
	GetID() string
	Action() string
}

type Outlet struct {
	Class string
	ID    string
}

// GetClass returns the class of the device
func (o Outlet) GetClass() string {
	return o.Class
}

// GetID returns the ID of the device
func (o Outlet) GetID() string {
	return o.ID
}

// Action represents an action performed by the outlet
func (o Outlet) Action() string {
	return "Powering on/off an appliance..."
}

// Light struct implementing the Device interface
type Light struct {
	Class string
	ID    string
}

// GetClass returns the class of the device
func (l Light) GetClass() string {
	return l.Class
}

// GetID returns the ID of the device
func (l Light) GetID() string {
	return l.ID
}

// Action represents an action performed by the light
func (l Light) Action() string {
	return "Turning light on/off..."
}

func NewDevice(class, id string) (Device, error) {
	switch class {
	case "Outlet":
		return Outlet{ID: id}, nil
	case "Light":
		return Light{ID: id}, nil
	default:
		return nil, errors.New("invalid device class")
	}
}
