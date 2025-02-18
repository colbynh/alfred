package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logrus logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	cfg := config{
		addr:     "0.0.0.0:8080",
		logLevel: "debug",
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	svr := app.mount()

	logger.Fatal(app.run(svr))
}
