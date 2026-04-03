// cmd/import imports Dublin cycle parking data from Smart Dublin's open dataset
// and inserts stands that don't already exist in the database.
//
// Usage:
//
//	go run ./cmd/import [--dry-run] [--radius 5] [--warn-radius 25] [--type "Sheffield Stand"] [--url <geojson-url>]
//
// Deduplication:
//   - Stands within --radius metres of an existing stand are skipped (true duplicates at the same rack).
//   - Stands within --warn-radius metres are imported but flagged — use this to catch
//     suspicious near-neighbours that may need manual review. Roads are typically 8–12m
//     wide, so opposite-side-of-road stands will often fall in this band.
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
	"math"
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
	radiusM := flag.Float64("radius", 5, "Skip import if an existing stand is within this many metres (same physical rack)")
	warnRadiusM := flag.Float64("warn-radius", 25, "Warn (but still import) if an existing stand is within this many metres")
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

	// Load all existing stands once so we can do in-process dedup checks.
	var existing []stand.Stand
	if err := db.Where("deleted_at IS NULL").Find(&existing).Error; err != nil {
		log.Fatalf("load existing stands: %v", err)
	}
	log.Printf("Loaded %d existing stands from database", len(existing))

	var (
		imported    int
		warned      int
		skippedDup  int
		skippedProx int
		skippedErr  int
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

		// 1. Skip if we already have a stand from this source with the same ID.
		if sourceID != "" && hasSourceID(existing, sourceName, sourceID) {
			skippedDup++
			continue
		}

		// 2. Skip if any existing stand is within the hard dedup radius (same physical rack).
		if near := nearestWithin(existing, lat, lng, *radiusM); near != nil {
			log.Printf("  SKIP  %.1fm → existing %s/%s: FID=%s lat=%.6f lng=%.6f",
				haversine(lat, lng, near.Lat, near.Lng), near.Source, near.SourceID, sourceID, lat, lng)
			skippedProx++
			continue
		}

		// 3. Warn (but still import) if a stand is suspiciously close — e.g. opposite side of road.
		if near := nearestWithin(existing, lat, lng, *warnRadiusM); near != nil {
			log.Printf("  WARN  %.1fm → existing %s/%s: FID=%s lat=%.6f lng=%.6f — importing anyway, check manually",
				haversine(lat, lng, near.Lat, near.Lng), near.Source, near.SourceID, sourceID, lat, lng)
			warned++
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
			log.Printf("  DRY-RUN insert: FID=%s lat=%.6f lng=%.6f name=%q stands=%d",
				sourceID, lat, lng, name, nStands)
		} else {
			if err := db.Create(&s).Error; err != nil {
				log.Printf("  ERROR inserting FID=%s: %v", sourceID, err)
				skippedErr++
				continue
			}
			// Add to in-memory slice so subsequent features can dedup against it.
			existing = append(existing, s)
		}
		imported++
	}

	fmt.Printf("\nDone.\n")
	fmt.Printf("  Imported:             %d\n", imported)
	fmt.Printf("  Warned (5–25m away):  %d\n", warned)
	fmt.Printf("  Skipped (dup ID):     %d\n", skippedDup)
	fmt.Printf("  Skipped (<5m away):   %d\n", skippedProx)
	fmt.Printf("  Skipped (error):      %d\n", skippedErr)
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

func hasSourceID(stands []stand.Stand, source, sourceID string) bool {
	for _, s := range stands {
		if s.Source == source && s.SourceID == sourceID {
			return true
		}
	}
	return false
}

func nearestWithin(stands []stand.Stand, lat, lng, radiusM float64) *stand.Stand {
	for i := range stands {
		if haversine(lat, lng, stands[i].Lat, stands[i].Lng) <= radiusM {
			return &stands[i]
		}
	}
	return nil
}

// haversine returns the distance in metres between two WGS-84 coordinates.
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const earthR = 6371000.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	return earthR * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
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
