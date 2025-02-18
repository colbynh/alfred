// Package outlet provides functionality for controlling smart outlets.
// It defines the core interface and factory methods for creating outlet controllers.
package outlet

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Outlet defines the interface for controlling smart outlets.
// Implementations must provide methods for basic device control,
// state management, and device discovery.
type Outlet interface {
	// getID returns the unique identifier of the outlet
	getID() string

	// getBrand returns the brand name of the outlet
	getBrand() string

	// action executes a command on the outlet and returns any error
	// Supported actions vary by implementation but typically include:
	// "on", "off", "state", "sysinfo", and "discover"
	action(action string, c *gin.Context) error

	// state retrieves the current state of the outlet
	// Returns a map containing at minimum a "state" field
	state() (map[string]interface{}, error)

	// sysInfo retrieves system information from the outlet
	// Returns a map containing device-specific details
	sysInfo() (map[string]interface{}, error)

	// discoverDevicesIps scans the network for compatible devices
	// Returns a map containing an "ips" array of discovered devices
	discoverDevicesIps() (map[string]interface{}, error)
}

// newOutlet creates a new Outlet instance based on the specified brand.
// Currently supported brands:
//   - "kasa": TP-Link Kasa smart outlets
//
// Parameters:
//   - brand: The outlet brand name (case-sensitive)
//   - id: Unique identifier for the device (typically IP address)
//   - c: Gin context for HTTP request handling
//   - logger: Logger for operation tracking
//
// Returns:
//   - Outlet: An implementation of the Outlet interface
//   - error: Non-nil if brand is unsupported
func newOutlet(brand string, id string, c *gin.Context, logger *logrus.Logger) (Outlet, error) {
	switch brand {
	case "kasa":
		return &kasaOutlet{id: id, c: c, logger: logger}, nil
	default:
		return nil, errors.New("unsupported outlet brand")
	}
}
