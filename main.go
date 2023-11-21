package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/semanticallynull/dublinbikeparking/apiv0"
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

	r := gin.New()

	apiRouter := r.Group("/api/v0")
	apiv0.NewAPIv0(apiRouter, db)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.Static("/static", StaticDirectoryV1)
	r.NoRoute(func(c *gin.Context) {
		c.File(StaticDirectoryV1 + "/index.html")
	})

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
