package outlet

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OutletActionHandler(svr *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		brand := c.Param("brand") // 'name' represents the brand
		action := c.Param("action")

		outlet, err := newOutlet(brand, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported outlet brand"})
			return
		}

		err = outlet.action(action)
		c.JSON(http.StatusOK, gin.H{
			"brand":  outlet.getBrand(),
			"id":     outlet.getID(),
			"action": action,
			"result": err,
		})
	}
}
