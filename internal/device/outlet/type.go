package outlet

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Outlet interface {
	getID() string
	getBrand() string
	action(action string, c *gin.Context) error
	state() (map[string]interface{}, error)
	sysInfo() (map[string]interface{}, error)
	discoverDevicesIps() (map[string]interface{}, error)
}

func newOutlet(brand string, id string, c *gin.Context, logger *logrus.Logger) (Outlet, error) {
	switch brand {
	case "kasa":
		return &kasaOutlet{id: id, c: c, logger: logger}, nil
	default:
		return nil, errors.New("unsupported outlet brand")
	}
}
