package apiv0

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	geojson "github.com/paulmach/go.geojson"
	validator "gopkg.in/validator.v2"
)

type Stand struct {
	gorm.Model
	Lat            float64 `validate:"min=-90,max=90"`
	Lng            float64 `validate:"min=-180,max=180"`
	Source         string
	SourceID       string
	Name           string
	Type           string `validate:"nonzero"`
	NumberOfStands int
	Notes          string
	Checked        string `json:"-"`
	Verified       bool
}

func (a *api) getStands(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	fc := geojson.NewFeatureCollection()

	dbc := a.DB.Where("checked != ?", "")

	if db := r.URL.Query().Get("dublinbikes"); db == "off" {
		dbc = dbc.Where("type != ?", "DublinBikes")
	} else if db == "only" {
		dbc = dbc.Where("type = ?", "DublinBikes")
	}

	stands := []Stand{}
	dbc.Find(&stands)

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
				"verified":       stand.Verified,
			},
		})
	}

	json.NewEncoder(w).Encode(fc)
}

func (a *api) createStand(w http.ResponseWriter, r *http.Request) {
	var stand Stand

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	json.Unmarshal(body, &stand)

	err = validator.Validate(&stand)
	if err != nil {
		errs := err.(validator.ErrorMap)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errs,
		})
		return
	}

	stand.Checked = ""
	stand.Source = "User Submission"
	stand.SourceID = uuid.New().String()[:7]

	a.DB.Create(&stand)

	json.NewEncoder(w).Encode(stand)
}
