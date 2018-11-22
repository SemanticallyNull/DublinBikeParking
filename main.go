package main

import (
	"log"
	"net/http"
	"os"

	"code.benchapman.ie/dublinbikeparking/apiv0"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var dialect = "sqlite3"
	var connectionString = "./demo.db"

	if dbDialect := os.Getenv("DBP_DB_DIALECT"); dbDialect != "" {
		dialect = dbDialect
	}
	if dbConnectionString := os.Getenv("DBP_DB_CONNECTION_STRING"); dbConnectionString != "" {
		connectionString = dbConnectionString
	}

	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api/v0").Subrouter()
	apiv0.NewAPIv0(apiRouter, db)

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
