package apiv0

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sendgrid/sendgrid-go/helpers/mail"

	validator "gopkg.in/validator.v2"

	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	geojson "github.com/paulmach/go.geojson"
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
	ImageID        string
	Notes          string
	Checked        string
	Verified       bool
	LastUpdateBy   string
	Thefts         []Theft
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

	dbc := a.DB

	if checked := r.URL.Query().Get("checked"); checked != "unchecked" {
		dbc = a.DB.Where("checked != ?", "")
	}

	if db := r.URL.Query().Get("dublinbikes"); db == "off" {
		dbc = dbc.Where("type != ?", "DublinBikes")
	} else if db == "only" {
		dbc = dbc.Where("type = ?", "DublinBikes")
	}

	if r.URL.Query().Get("review") == "true" {
		dbc = dbc.Where("number_of_stands IS NULL")
	}

	stands := []Stand{}
	dbc.Preload("Thefts").Find(&stands)

	for _, stand := range stands {
		feature := &geojson.Feature{
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
				"imageId":        stand.ImageID,
				"source":         stand.Source,
				"checked":        stand.Checked != "",
				"verified":       stand.Verified,
				"thefts":         []map[string]interface{}{},
			},
		}

		for _, theft := range stand.Thefts {
			feature.Properties["thefts"] = append(feature.Properties["thefts"].([]map[string]interface{}), map[string]interface{}{
				"id": theft.ID,
			})
		}

		fc.AddFeature(feature)
	}

	json.NewEncoder(w).Encode(fc)
}

func (a *api) getStand(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	vars := mux.Vars(r)

	stand := &Stand{}
	a.DB.Where("`stand_id` = ?", vars["id"]).Preload("Thefts").First(stand)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":             stand.StandID,
		"name":           stand.Name,
		"type":           stand.Type,
		"numberOfStands": stand.NumberOfStands,
		"notes":          stand.Notes,
		"source":         stand.Source,
		"checked":        stand.Checked != "",
		"verified":       stand.Verified,
		"thefts":         stand.Thefts,
	})
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

	go func() {
		from := mail.NewEmail("DublinBikeParking", "no-reply@dublinbikeparking.com")
		subject := fmt.Sprintf("New Stand at '%s'", stand.Name)
		to := mail.NewEmail("Katie Chapman", "hello@katiechapman.ie")
		plainTextContent := "New stand on DublinBikeParking.com\n" +
			"Stand ID: " + stand.StandID + "\n" +
			"Name: " + stand.Name + "\n" +
			"Coordinates: " + fmt.Sprintf("%f %f", stand.Lat, stand.Lng) + "\n" +
			"https://dublinbikeparking.com/update.html#19/" + fmt.Sprintf("%f/%f", stand.Lat, stand.Lng)
		htmlContent := "New stand on DublinBikeParking.com<br>" +
			"<b>Stand ID:</b> " + stand.StandID + "<br>" +
			"<b>Name:</b> " + stand.Name + "<br>" +
			"<b>Coordinates:</b> " + fmt.Sprintf("%f %f", stand.Lat, stand.Lng) + "<br>" +
			"<a href=\"https://dublinbikeparking.com/update.html#19/" + fmt.Sprintf("%f/%f", stand.Lat, stand.Lng) + "\">Link to update page</a>"
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(a.SendgridAPIKey)
		_, err := client.Send(message)
		if err != nil {
			log.Println(err)
		}
	}()

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

	userEmail := context.Get(r, "userEmail").(string)
	sub := context.Get(r, "userSub").(string)

	checked := ""

	if updatedStand.Checked != "" {
		checked = sub
	}

	if originStand.ID != 0 {

		tx := a.DB.Begin()
		err = tx.Model(&originStand).Limit(1).Update(map[string]interface{}{
			"name":             updatedStand.Name,
			"type":             updatedStand.Type,
			"number_of_stands": updatedStand.NumberOfStands,
			"checked":          checked,
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

func (a *api) deleteStand(w http.ResponseWriter, r *http.Request) {
	var stand Stand

	vars := mux.Vars(r)

	a.DB.Where("stand_id = ?", vars["id"]).First(&stand)

	a.DB.Delete(&stand)

	w.WriteHeader(http.StatusAccepted)
}
