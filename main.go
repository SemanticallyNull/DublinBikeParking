package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/packr"
)

func main() {
	fs := http.FileServer(packr.NewBox("./static"))
	http.Handle("/", fs)

	port := getPort()
	log.Printf("Listening on port %s...", port)
	http.ListenAndServe(":"+port, nil)
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "3000"
}
