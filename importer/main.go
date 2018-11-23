package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"code.benchapman.ie/dublinbikeparking/apiv0"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	db.LogMode(true)

	r, _ := os.Open(os.Args[1])
	records, _ := csv.NewReader(r).ReadAll()
	for _, record := range records {
		var stands []apiv0.Stand

		db.Where("lat LIKE ? AND lng LIKE ?", record[1][:8]+"%", record[0][:8]+"%").Find(&stands)

		noStands, err := strconv.Atoi(strings.TrimSuffix(record[7], ".0"))
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		if len(stands) == 0 {
			lat, _ := strconv.ParseFloat(record[1], 64)
			lng, _ := strconv.ParseFloat(record[0], 64)

			db.Create(&apiv0.Stand{
				Model: gorm.Model{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Lat:            lat,
				Lng:            lng,
				Name:           record[2],
				SourceID:       record[4],
				NumberOfStands: noStands,
				Type:           record[8],
			})
		}
	}
}
