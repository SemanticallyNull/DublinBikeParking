package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	port := getPort()
	log.Printf("Listening on port %s...", port)
	http.ListenAndServe(":"+port, r)
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "3000"
}
