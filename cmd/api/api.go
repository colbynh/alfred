package main

import (
	// "github.com/colbynh/alfred/internal/device/light"
	"github.com/colbynh/alfred/internal/device/outlet"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	// "net/http"
)

type application struct {
	config config
	logger *logrus.Logger
}

type config struct {
	addr     string
	env      string
	apiURL   string
	logLevel string
}

func (app *application) mount() *gin.Engine {
	svr := gin.New()
	svr.Use(logger.SetLogger())

	// Middleware to require a header
	// authHeader := func(c *gin.Context) {
	// 	if c.GetHeader("hue-application-key") == "" {
	// 		app.logger.Warn("Missing hue-application-key header")
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "hue-application-key is required"})
	// 		c.Abort()
	// 		return
	// 	}
	// 	c.Next()
	// }

	// Apply the middleware to specific routes
	svr.POST("/api/v1/device/outlet/:brand/:id/:action", outlet.OutletActionHandler(svr, app.logger))
	svr.GET("/api/v1/device/outlet/:brand/:id/:action", outlet.OutletActionHandler(svr, app.logger))
	// svr.GET("/api/v1/device/outlet/:brand/:action", outlet.OutletActionHandler(svr))

	// svr.PUT("/api/v1/device/light/:brand/:ip/:id/:action", authHeader, light.LightActionHandler(svr, app.logger))
	// svr.GET("/api/v1/device/light/:brand/:ip/:action", authHeader, light.LightActionHandler(svr, logger))
	// TODO: add delete route and test

	return svr
}

func (app *application) run(svr *gin.Engine) error {
	app.logger.Info("Starting server on", app.config.addr)
	if err := svr.Run(app.config.addr); err != nil {
		app.logger.Error("Error starting server:", err)
		return err
	}
	app.logger.Info("Server stopped")
	return nil
}
