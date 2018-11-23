package apiv0

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	geojson "github.com/paulmach/go.geojson"
)

type Stand struct {
	gorm.Model
	Lat            float64
	Lng            float64
	Source         string
	SourceID       string
	Name           string
	Type           string
	NumberOfStands int
	Notes          string
}

func (a *api) getStands(w http.ResponseWriter, r *http.Request) {
	fc := geojson.NewFeatureCollection()

	stands := []Stand{}
	a.DB.Find(&stands)

	for _, stand := range stands {
		fc.AddFeature(&geojson.Feature{
			Geometry: &geojson.Geometry{
				Type:  geojson.GeometryPoint,
				Point: []float64{stand.Lng, stand.Lat},
			},
			Properties: map[string]interface{}{
				"id":             stand.ID,
				"name":           stand.Name,
				"type":           stand.Type,
				"numberOfStands": stand.NumberOfStands,
				"notes":          stand.Notes,
				"source":         stand.Source,
			},
		})
	}

	json.NewEncoder(w).Encode(fc)
}
