package apiv0

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	geojson "github.com/paulmach/go.geojson"
	validator "gopkg.in/validator.v2"
)

type Stand struct {
	gorm.Model
	StandID        string
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
	LastUpdateBy   string
}

type StandUpdate struct {
	gorm.Model
	Stand     Stand
	StandID   uint
	UserEmail string
	Update    string `sql:"type:text"`
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

	if r.URL.Query().Get("review") == "true" {
		dbc = dbc.Where("number_of_stands IS NULL")
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
				"id":             stand.StandID,
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
	id := uuid.New().String()[:6]
	stand.SourceID = id
	stand.StandID = id

	a.DB.Create(&stand)

	json.NewEncoder(w).Encode(stand)
}

func (a *api) updateStand(w http.ResponseWriter, r *http.Request) {
	var originStand, updatedStand Stand

	vars := mux.Vars(r)

	a.DB.Where("stand_id = ?", vars["id"]).First(&originStand)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.Unmarshal(body, &updatedStand)

	err = validator.Validate(&updatedStand)
	if err != nil {
		errs := err.(validator.ErrorMap)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errs,
		})
		return
	}

	if originStand.ID != 0 {
		userEmail := context.Get(r, "userEmail").(string)

		tx := a.DB.Begin()
		err = tx.Model(&originStand).Limit(1).Update(map[string]interface{}{
			"name":             updatedStand.Name,
			"type":             updatedStand.Type,
			"number_of_stands": updatedStand.NumberOfStands,
			"last_update_by":   userEmail,
		}).Error
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		updateJson, err := json.Marshal(map[string]interface{}{
			"name":             updatedStand.Name,
			"type":             updatedStand.Type,
			"number_of_stands": updatedStand.NumberOfStands,
		})
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		update := StandUpdate{
			Stand:     originStand,
			UserEmail: userEmail,
			Update:    string(updateJson),
		}
		if err := tx.Create(&update).Error; err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		tx.Commit()
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []string{
				"stand not found",
			},
		})
	}

	json.NewEncoder(w).Encode(originStand)
}
