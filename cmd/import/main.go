// cmd/import imports Dublin cycle parking data from Smart Dublin's open dataset
// and inserts stands that don't already exist in the database.
//
// Usage:
//
//	go run ./cmd/import [--dry-run] [--type "Sheffield Stand"] [--url <geojson-url>]
//
// Deduplication is by source ID only (source=smartdublin + FID). This makes
// the import safe to re-run. Proximity-based dedup is intentionally omitted —
// Dublin has narrow 1-way streets and opposite-side-of-road stands can be less
// than 3m apart, making any distance threshold unreliable. Use --dry-run first
// and review any spatial duplicates on the map afterwards.
//
// Environment variables (same as main server):
//
//	DBP_DB_DIALECT          sqlite3 (default) | mysql
//	DBP_DB_CONNECTION_STRING ./demo.db (default)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/semanticallynull/dublinbikeparking/stand"
)

const defaultURL = "https://data.smartdublin.ie/dataset/eb3c5fcb-5df4-4993-bf3a-c07afea32397/resource/95757080-3adc-425c-bdd1-73c9e6c33fc2/download/dublin-public-cycle-parking-facilities.geojson"
const sourceName = "smartdublin"

// GeoJSON types — just enough to parse the Smart Dublin feed.
type featureCollection struct {
	Features []feature `json:"features"`
}

type feature struct {
	Properties map[string]interface{} `json:"properties"`
	Geometry   geometry               `json:"geometry"`
}

type geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"` // [lng, lat]
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Print what would be inserted without writing to the database")
	defaultType := flag.String("type", "Sheffield Stand", "Stand type to assign when the source data doesn't specify one")
	sourceURL := flag.String("url", defaultURL, "GeoJSON source URL")
	flag.Parse()

	dialect := "sqlite3"
	connStr := "./demo.db"
	if v := os.Getenv("DBP_DB_DIALECT"); v != "" {
		dialect = v
	}
	if v := os.Getenv("DBP_DB_CONNECTION_STRING"); v != "" {
		connStr = v
	}

	db, err := gorm.Open(dialect, connStr)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	log.Printf("Fetching %s", *sourceURL)
	fc, err := fetchGeoJSON(*sourceURL)
	if err != nil {
		log.Fatalf("fetch: %v", err)
	}
	log.Printf("Fetched %d features", len(fc.Features))

	// Build a set of already-imported source IDs for O(1) lookup.
	imported_ids := map[string]bool{}
	var existing []stand.Stand
	if err := db.Where("source = ? AND deleted_at IS NULL", sourceName).Find(&existing).Error; err != nil {
		log.Fatalf("load existing stands: %v", err)
	}
	for _, s := range existing {
		imported_ids[s.SourceID] = true
	}
	log.Printf("Found %d existing stands from source=%s", len(existing), sourceName)

	var (
		imported   int
		skippedDup int
		skippedErr int
	)

	for _, f := range fc.Features {
		if f.Geometry.Type != "Point" || len(f.Geometry.Coordinates) < 2 {
			skippedErr++
			continue
		}

		lng := f.Geometry.Coordinates[0]
		lat := f.Geometry.Coordinates[1]
		sourceID := propString(f.Properties, "FID")
		name := propString(f.Properties, "Location")
		nStands := propInt(f.Properties, "nostands")

		// Skip if already imported (safe re-runs).
		if sourceID != "" && imported_ids[sourceID] {
			skippedDup++
			continue
		}

		s := stand.Stand{
			StandID:        uuid.New().String(),
			Lat:            lat,
			Lng:            lng,
			Source:         sourceName,
			SourceID:       sourceID,
			Name:           name,
			Type:           *defaultType,
			NumberOfStands: nStands,
		}

		if *dryRun {
			log.Printf("  INSERT FID=%-6s  stands=%-3d  lat=%.6f  lng=%.6f  name=%q",
				sourceID, nStands, lat, lng, name)
		} else {
			if err := db.Create(&s).Error; err != nil {
				log.Printf("  ERROR inserting FID=%s: %v", sourceID, err)
				skippedErr++
				continue
			}
			imported_ids[sourceID] = true
		}
		imported++
	}

	fmt.Printf("\nDone.\n")
	fmt.Printf("  Imported:          %d\n", imported)
	fmt.Printf("  Skipped (dup ID):  %d\n", skippedDup)
	fmt.Printf("  Skipped (error):   %d\n", skippedErr)
	if *dryRun {
		fmt.Println("  (dry-run — nothing written)")
	}
}

func fetchGeoJSON(url string) (*featureCollection, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var fc featureCollection
	if err := json.NewDecoder(resp.Body).Decode(&fc); err != nil {
		return nil, err
	}
	return &fc, nil
}

func propString(props map[string]interface{}, key string) string {
	v, ok := props[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func propInt(props map[string]interface{}, key string) int {
	v, ok := props[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case string:
		i, _ := strconv.Atoi(n)
		return i
	}
	return 0
}
