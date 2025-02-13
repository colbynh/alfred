package light

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type light interface {
	on() error
	off() error
	setBrightness() error
	color() error
	execAction(action string) error
	getGinContext() *gin.Context
}

func newLight(brand string, ip string, id string, ctx *gin.Context) (light, error) {
	switch brand {
	case "philips":
		return &philipsLight{brand: brand, ip: ip, id: id, ctx: ctx}, nil
	default:
		return nil, errors.New("unsupported light brand")
	}
}
