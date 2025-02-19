// Package outlet provides functionality for controlling smart outlets.
// It supports different brands of smart outlets through a common interface
// and provides HTTP handlers for device control.
package outlet

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// OutletActionHandler creates a gin.HandlerFunc that processes outlet control requests.
// It handles device actions like power on/off, state queries, and device discovery.
//
// Parameters:
//   - svr: The gin engine instance for HTTP routing
//   - logger: A configured logrus logger for operation tracking
//
// The handler expects URL parameters:
//   - brand: The outlet brand (e.g., "kasa")
//   - id: Device identifier (typically IP address)
//   - action: Command to execute (e.g., "on", "off", "state")
//
// Returns a gin.HandlerFunc that:
//   - Creates an appropriate outlet controller based on brand
//   - Executes the requested action
//   - Returns JSON response with operation result
//
// Example URL: POST /api/v1/device/outlet/kasa/192.168.1.100/on
func OutletActionHandler(svr *gin.Engine, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		brand := c.Param("brand")
		action := c.Param("action")

		logger.Debugf("Received request: brand=%s, 'id=%s', 'action=%s'", brand, id, action)
		outlet, err := newOutlet(brand, id, c, logger)
		if err != nil {
			logger.Errorf("Error creating outlet: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported outlet brand"})
			return
		}

		err = outlet.action(action, c)
		if err != nil {
			logger.Errorf("Error executing action: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
}
