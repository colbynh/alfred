package main

import (
	"github.com/colbynh/alfred/internal/device/outlet"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

type application struct {
	config config
	// store  store.Storage
	// logger *zap.SugaredLogger
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *gin.Engine {
	svr := gin.New()
	svr.Use(logger.SetLogger())

	svr.POST("/devices/outlets/:id/:brand/:action", OutletActionHandler(svr))
}

func (app *application) run() error {
	svr := gin.New()
	svr.Use(logger.SetLogger())

	// Listen and Server in 0.0.0.0:8080
	if err := svr.Run(app.config.addr); err != nil {
		return err
	}
	return nil
}
