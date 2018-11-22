package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"code.benchapman.ie/dublinbikeparking/apiv0"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var dialect = "sqlite3"
	var connectionString = "./demo.db"

	if cfenv.IsRunningOnCF() {
		cfapp, err := cfenv.Current()
		if err != nil {
			panic(err)
		}

		service, err := cfapp.Services.WithName("dublinbikeparking-db")
		if err != nil {
			log.Fatalf("%s", err)
		}

		dialect = "mysql"
		hostname, _ := service.CredentialString("hostname")
		username, _ := service.CredentialString("username")
		password, _ := service.CredentialString("password")
		dbname, _ := service.CredentialString("name")
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", username, password, hostname, dbname)
	}

	if dbDialect := os.Getenv("DBP_DB_DIALECT"); dbDialect != "" {
		connectionString = dbDialect
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
