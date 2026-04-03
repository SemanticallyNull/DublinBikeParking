// cmd/import imports Dublin cycle parking data from Smart Dublin's open dataset
// and inserts stands that don't already exist in the database.
//
// Usage:
//
//	go run ./cmd/import [--dry-run] [--type "Sheffield Stand"] [--url <geojson-url>]
//
// Deduplication uses two signals combined:
//
//	Distance alone    ≤5m              → skip  (same physical rack regardless of count)
//	Distance + count  ≤20m, same count → skip  (same stand, counted consistently)
//	Distance alone    ≤20m             → warn  (close but counts differ — opposite side of road?)
//	Distance + count  ≤50m, same count → warn  (same count, suspiciously nearby)
//
// Stand count is only used as a signal when both sides have a recorded value > 0.
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
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/semanticallynull/dublinbikeparking/stand"
)

const defaultURL = "https://data.smartdublin.ie/dataset/eb3c5fcb-5df4-4993-bf3a-c07afea32397/resource/95757080-3adc-425c-bdd1-73c9e6c33fc2/download/dublin-public-cycle-parking-facilities.geojson"
const sourceName = "smartdublin"

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

	var existing []stand.Stand
	if err := db.Where("deleted_at IS NULL").Find(&existing).Error; err != nil {
		log.Fatalf("load existing stands: %v", err)
	}
	log.Printf("Loaded %d existing stands from database", len(existing))

	// Index existing source IDs for O(1) exact-match lookup.
	existingSourceIDs := map[string]bool{}
	for _, s := range existing {
		if s.Source == sourceName {
			existingSourceIDs[s.SourceID] = true
		}
	}

	type warningEntry struct {
		dist     float64
		msg      string
		sourceID string
		lat, lng float64
		nStands  int
		name     string
	}
	var warnings []warningEntry

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

		// 1. Exact source ID match — already imported, skip.
		if sourceID != "" && existingSourceIDs[sourceID] {
			skippedDup++
			continue
		}

		// 2. Check against all existing stands using combined signals.
		skip, warnMsg := checkDuplicate(lat, lng, nStands, existing)
		if skip {
			log.Printf("  SKIP  %s  FID=%s lat=%.6f lng=%.6f", warnMsg, sourceID, lat, lng)
			skippedDup++
			continue
		}
		if warnMsg != "" {
			dist := nearestDist(lat, lng, existing)
			warnings = append(warnings, warningEntry{dist, warnMsg, sourceID, lat, lng, nStands, name})
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
			existing = append(existing, s)
			existingSourceIDs[sourceID] = true
		}
		imported++
	}

	// Write warnings to file sorted by distance (closest first).
	if len(warnings) > 0 {
		sort.Slice(warnings, func(i, j int) bool { return warnings[i].dist < warnings[j].dist })
		wf, err := os.Create("import-warnings.txt")
		if err != nil {
			log.Printf("WARNING: could not create import-warnings.txt: %v", err)
		} else {
			fmt.Fprintf(wf, "%-8s  %-6s  %-4s  %-10s  %-10s  %-30s  %s\n",
				"DIST(m)", "FID", "CNT", "LAT", "LNG", "NAME", "REASON")
			fmt.Fprintf(wf, "%s\n", strings.Repeat("-", 100))
			for _, w := range warnings {
				fmt.Fprintf(wf, "%-8.1f  %-6s  %-4d  %-10.6f  %-10.6f  %-30s  %s\n",
					w.dist, w.sourceID, w.nStands, w.lat, w.lng, w.name, w.msg)
			}
			wf.Close()
			log.Printf("Wrote %d warnings to import-warnings.txt", len(warnings))
		}
	}

	fmt.Printf("\nDone.\n")
	fmt.Printf("  Imported:         %d\n", imported)
	fmt.Printf("  Warned:           %d\n", len(warnings))
	fmt.Printf("  Skipped (dup):    %d\n", skippedDup)
	fmt.Printf("  Skipped (error):  %d\n", skippedErr)
	if *dryRun {
		fmt.Println("  (dry-run — nothing written)")
	}
}

func nearestDist(lat, lng float64, existing []stand.Stand) float64 {
	min := math.MaxFloat64
	for i := range existing {
		if d := haversine(lat, lng, existing[i].Lat, existing[i].Lng); d < min {
			min = d
		}
	}
	return min
}

// checkDuplicate returns (skip, description) by combining distance and stand count.
// Stand count is only used as a signal when both sides have a recorded value > 0.
func checkDuplicate(lat, lng float64, nStands int, existing []stand.Stand) (skip bool, msg string) {
	for i := range existing {
		e := &existing[i]
		dist := haversine(lat, lng, e.Lat, e.Lng)
		sameCount := nStands > 0 && e.NumberOfStands > 0 && nStands == e.NumberOfStands

		switch {
		case dist <= 5:
			// Within 5m — same physical rack regardless of count.
			return true, fmt.Sprintf("%.1fm from %s/%s (too close)", dist, e.Source, e.SourceID)

		case dist <= 20 && sameCount:
			// Close + identical stand count — almost certainly the same stand.
			return true, fmt.Sprintf("%.1fm from %s/%s, same count (%d)", dist, e.Source, e.SourceID, nStands)

		case dist <= 20:
			// Close but counts differ — could be opposite side of road.
			return false, fmt.Sprintf("%.1fm from %s/%s (counts differ: incoming %d, existing %d)", dist, e.Source, e.SourceID, nStands, e.NumberOfStands)

		case dist <= 50 && sameCount:
			// Same count, suspiciously nearby.
			return false, fmt.Sprintf("%.1fm from %s/%s, same count (%d)", dist, e.Source, e.SourceID, nStands)
		}
	}
	return false, ""
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
