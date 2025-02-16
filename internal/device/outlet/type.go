package outlet

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Outlet interface {
	getID() string
	getBrand() string
	action(action string, c *gin.Context) error
	state() (map[string]interface{}, error)
	sysInfo() (map[string]interface{}, error)
}

func newOutlet(brand string, id string, c *gin.Context) (Outlet, error) {
	switch brand {
	case "kasa":
		return &kasaOutlet{id: id, c: c}, nil
	default:
		return nil, errors.New("unsupported outlet brand")
	}
}
