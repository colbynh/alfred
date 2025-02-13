package light

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LightActionHandler(svr *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		brand := c.Param("brand")
		ip := c.Param("ip")
		action := c.Param("action")
		jsonPayload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		light, err := newLight(brand, ip, id, action, string(jsonPayload))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported light brand"})
			return
		}
		err = light.execAction(action)
		c.JSON(http.StatusOK, gin.H{
			"brand":  brand,
			"ip":     ip,
			"id":     id,
			"action": action,
			"body":   jsonPayload,
			"result": err,
		})
	}
}
