package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

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

	r, _ := os.Open("import.csv")
	records, _ := csv.NewReader(r).ReadAll()
	for ln, record := range records {
		var stand apiv0.Stand

		fmt.Printf("%d - %#v", ln, record)
		lat := record[27][1:9]
		lng := record[28][:8]
		db.Where("lat LIKE ? AND lng LIKE ?", lat+"%", lng+"%").First(&stand)

		stands, _ := strconv.Atoi(record[6])
		db.Model(&stand).Where("lat LIKE ? AND lng LIKE ?", lat+"%", lng+"%").Updates(map[string]interface{}{
			"name":             record[5],
			"number_of_stands": stands,
			"type":             record[0],
			"source_id":        fmt.Sprintf("%d", ln),
		})

		//		noStands, err := strconv.Atoi(strings.TrimSuffix(record[7], ".0"))
		//		if err != nil {
		//			fmt.Printf("%s\n", err)
		//		}
		//
		//		if len(stands) == 0 {
		//			lat, _ := strconv.ParseFloat(record[1], 64)
		//			lng, _ := strconv.ParseFloat(record[0], 64)
		//
		//			db.Create(&apiv0.Stand{
		//				Model: gorm.Model{
		//					CreatedAt: time.Now(),
		//					UpdatedAt: time.Now(),
		//				},
		//				Lat:            lat,
		//				Lng:            lng,
		//				Name:           record[2],
		//				SourceID:       record[4],
		//				NumberOfStands: noStands,
		//				Type:           record[8],
		//			})
		//		}
	}
}
