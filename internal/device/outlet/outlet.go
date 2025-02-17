package outlet

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func OutletActionHandler(svr *gin.Engine, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		brand := c.Param("brand") // 'name' represents the brand
		action := c.Param("action")

		logger.Debugf("Received request: brand=%s, id=%s, action=%s", brand, id, action)

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

		logger.Debugf("Successfully executed action: %s for outlet: %s", action, id)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
