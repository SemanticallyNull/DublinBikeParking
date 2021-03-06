package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"code.katiechapman.ie/dublinbikeparking/apiv0"
)

const StaticDirectoryV1 = "./static"

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
		_, err := w.Write([]byte("ok"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("could not write: %s", err)
			return
		}
	})

	fs := http.FileServer(http.Dir(StaticDirectoryV1))
	r.PathPrefix("/").Handler(fs)

	port := getPort()
	log.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "3000"
}
