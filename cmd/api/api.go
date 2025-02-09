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

	svr.POST("/device/outlet/:brand/:id/:action", outlet.OutletActionHandler(svr))
	return svr
}

func (app *application) run(svr *gin.Engine) error {
	svr.Use(logger.SetLogger())
	if err := svr.Run(app.config.addr); err != nil {
		return err
	}
	return nil
}
