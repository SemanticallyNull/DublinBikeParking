package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"code.katiechapman.ie/dublinbikeparking/apiv0"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const StaticDirectoryV1 = "./static"
const StaticDirectoryV2 = "./static-vue/dist"

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

	if os.Getenv("DBP_UI_V2") == "true" {
		r.NotFoundHandler = r.NewRoute().HandlerFunc(serverHandler).GetHandler()
	} else {
		fs := http.FileServer(http.Dir(StaticDirectoryV1))
		r.PathPrefix("/").Handler(fs)
	}

	port := getPort()
	log.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(":"+port, handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"https://dublinbikeparking.com/", "https://dbpimg.apps.katiechapman.ie/"}),
	)(r))
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(StaticDirectoryV2 + r.URL.Path); err != nil {
		http.ServeFile(w, r, StaticDirectoryV2+"/index.html")
		return
	}
	http.ServeFile(w, r, StaticDirectoryV2+r.URL.Path)
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "3000"
}
