package outlet

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// OutletActionHandler processes API calls for outlets
func OutletActionHandler(svr *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(err)
		}
		brand := c.Param("brand") // 'name' represents the brand
		action := c.Param("action")

		// Create the correct outlet based on brand
		outlet, err := newOutlet(brand, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported outlet brand"})
			return
		}

		// Execute the action
		err = outlet.action(action)
		c.JSON(http.StatusOK, gin.H{
			"brand":  outlet.getBrand(),
			"id":     outlet.getID(),
			"action": action,
			"result": err,
		})
	}
}
