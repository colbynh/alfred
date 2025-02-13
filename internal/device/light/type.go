package light

import (
	"errors"
)

type light interface {
	on() error
	off() error
	setBrightness() error
	color() error
	execAction(action string) error
}

func newLight(brand, ip, id, action, jsonPayload string) (light, error) {
	switch brand {
	case "philips":
		return &philipsLight{brand: brand, ip: ip, id: id, jsonPayload: jsonPayload}, nil
	default:
		return nil, errors.New("unsupported light brand")
	}
}
