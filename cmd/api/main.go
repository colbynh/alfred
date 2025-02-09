package main

import (
	"log"
)

func main() {
	cfg := config{
		addr: "0.0.0.0:8080",
	}

	app := &application{
		config: cfg,
	}

	svr := app.mount()

	log.Fatal(app.run(svr))
}
