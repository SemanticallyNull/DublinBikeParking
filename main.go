package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	r.Use(func(c *gin.Context) {
		if c.GetHeader("X-Forwarded-Proto") == "http" {
			target := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, target)
			c.Abort()
			return
		}
		c.Next()
	})

	apiRouter := r.Group("/api/v0")
	apiv0.NewAPIv0(apiRouter, db)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.Static("/static", StaticDirectoryV1)
	r.Static("/assets", StaticDirectoryV1+"/assets")
	r.Static("/icons", StaticDirectoryV1+"/icons")
	r.StaticFile("/og-image.png", StaticDirectoryV1+"/og-image.png")
	r.StaticFile("/manifest.webmanifest", StaticDirectoryV1+"/manifest.webmanifest")
	r.StaticFile("/sw.js", StaticDirectoryV1+"/sw.js")
	r.StaticFile("/registerSW.js", StaticDirectoryV1+"/registerSW.js")
	r.NoRoute(func(c *gin.Context) {
		// Serve real files (e.g. workbox-*.js, hashed assets) if they exist,
		// otherwise fall back to index.html for SPA routing.
		filePath := filepath.Join(StaticDirectoryV1, filepath.Clean(c.Request.URL.Path))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			c.File(filePath)
			return
		}
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
