package main

import (
	"fmt"
	"log"
	"net/http"
)

const wedPort = "80"

type Config struct {
}

func main() {
	app := Config{}
	log.Printf("Starting broker service on port %s\n", wedPort)

	// Define http server
	src := &http.Server{
		Addr:    fmt.Sprintf(":%s", wedPort),
		Handler: app.routes(),
	}

	// Start the server
	err := src.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
