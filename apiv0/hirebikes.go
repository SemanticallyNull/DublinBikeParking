package apiv0

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	geojson "github.com/paulmach/go.geojson"
)

type dbStand struct {
	Name            string `json:"address"`
	Position        dbPos  `json:"position"`
	NoStands        int    `json:"no_stands"`
	AvailableStands int    `json:"available_bike_stands"`
	AvailableBikes  int    `json:"available_bikes"`
}
type dbPos struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type syncFeatureCollection struct {
	*geojson.FeatureCollection
	mu sync.Mutex
}

func newSyncFeatureCollection() syncFeatureCollection {
	return syncFeatureCollection{
		FeatureCollection: geojson.NewFeatureCollection(),
		mu:                sync.Mutex{},
	}
}

func (fc *syncFeatureCollection) AddFeature(feature *geojson.Feature) *syncFeatureCollection {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.FeatureCollection = fc.FeatureCollection.AddFeature(feature)
	return fc
}
func (fc *syncFeatureCollection) MarshalJSON() ([]byte, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	return fc.FeatureCollection.MarshalJSON()
}

func (a *api) getHireBikes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Add("content-type", "application/json")
	fc := newSyncFeatureCollection()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		resp, err := http.Get("https://data.smartdublin.ie/bleeperbike-api/last_snapshot/")
		if err != nil {
			fmt.Println(err)
			return err
		}

		var data []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err
		}

		for _, bike := range data {
			lng := bike["lon"].(float64)
			lat := bike["lat"].(float64)
			if lat == 0 || lng == 0 {
				continue
			}
			feature := &geojson.Feature{
				Geometry: &geojson.Geometry{
					Type:  geojson.GeometryPoint,
					Point: []float64{lng, lat},
				},
				Properties: map[string]interface{}{
					"type":     "Bleeper",
					"verified": true,
				},
			}
			fc.AddFeature(feature)
		}

		return err
	})
	//g.Go(func() error {
	//	resp, err := http.Get("https://data.smartdublin.ie/mobybikes-api/last_reading/")
	//	if err != nil && resp.StatusCode == 200 {
	//		return err
	//	}
	//
	//	var data []map[string]interface{}
	//	err = json.NewDecoder(resp.Body).Decode(&data)
	//	if err != nil {
	//		return err
	//	}
	//
	//	for _, bike := range data {
	//		lng := bike["Longitude"].(float64)
	//		lat := bike["Latitude"].(float64)
	//		if lat == 0 || lng == 0 {
	//			continue
	//		}
	//		feature := &geojson.Feature{
	//			Geometry: &geojson.Geometry{
	//				Type:  geojson.GeometryPoint,
	//				Point: []float64{lng, lat},
	//			},
	//			Properties: map[string]interface{}{
	//				"type":     "Moby",
	//				"verified": true,
	//			},
	//		}
	//		fc.AddFeature(feature)
	//	}
	//
	//	return err
	//})
	g.Go(func() error {
		dbAPIKey := os.Getenv("DUBLINBIKES_API_KEY")
		if dbAPIKey == "" {
			return nil
		}

		resp, err := http.Get(fmt.Sprintf("https://api.jcdecaux.com/vls/v1/stations?contract=dublin&apiKey=%s", dbAPIKey))
		if err != nil && resp.StatusCode == 200 {
			return err
		}

		var data []dbStand
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err
		}

		for _, bike := range data {
			feature := &geojson.Feature{
				Geometry: &geojson.Geometry{
					Type:  geojson.GeometryPoint,
					Point: []float64{bike.Position.Lng, bike.Position.Lat},
				},
				Properties: map[string]interface{}{
					"name":            bike.Name,
					"type":            "DublinBikes",
					"numberOfStands":  bike.NoStands,
					"bikesAvailable":  bike.AvailableBikes,
					"standsAvailable": bike.AvailableStands,
					"verified":        true,
				},
			}
			fc.AddFeature(feature)
		}

		return err
	})

	err := g.Wait()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			fmt.Printf("error writing error: %s", err)
			return
		}
		return
	}

	err = json.NewEncoder(w).Encode(fc)
	if err != nil {
		fmt.Printf("error encoding json: %s", err)
	}
}
