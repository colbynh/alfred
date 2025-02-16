package main

import (
	"github.com/colbynh/alfred/internal/device/light"
	"github.com/colbynh/alfred/internal/device/outlet"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type application struct {
	config config
	// store  store.Storage
	// logger *zap.SugaredLogger
}

type config struct {
	addr   string
	env    string
	apiURL string
}

func (app *application) mount() *gin.Engine {
	svr := gin.New()
	svr.Use(logger.SetLogger())

	// Middleware to require a header
	authHeader := func(c *gin.Context) {
		if c.GetHeader("hue-application-key") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Required-Header is required"})
			c.Abort()
			return
		}
		c.Next()
	}

	// Apply the middleware to specific routes
	svr.POST("/api/v1/device/outlet/:brand/:id/:action", outlet.OutletActionHandler(svr))
	svr.GET("/api/v1/device/outlet/:brand/:id/:action", outlet.OutletActionHandler(svr))

	svr.PUT("/api/v1/device/light/:brand/:ip/:id/:action", authHeader, light.LightActionHandler(svr))
	svr.GET("/api/v1/device/light/:brand/:ip/:action", authHeader, light.LightActionHandler(svr))
	// TODO: add delete route and test

	return svr
}

func (app *application) run(svr *gin.Engine) error {
	svr.Use(logger.SetLogger())
	if err := svr.Run(app.config.addr); err != nil {
		return err
	}
	return nil
}
