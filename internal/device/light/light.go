package light

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LightActionHandler(svr *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		brand := c.Param("brand")
		ip := c.Param("ip")
		action := c.Param("action")
		fmt.Println("debugging working!!!!!!!!!!!!!")

		light, err := newLight(brand, ip, id, c)
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
			"body":   c.Request.Body,
			"result": err,
		})
	}
}
