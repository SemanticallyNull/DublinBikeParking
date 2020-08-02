package apiv0

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	geojson "github.com/paulmach/go.geojson"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/validator.v2"

	"code.katiechapman.ie/dublinbikeparking/stand"
)

type StandUpdate struct {
	gorm.Model
	Stand     stand.Stand
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

	stands := []stand.Stand{}
	dbc.Find(&stands)

	for _, s := range stands {
		feature := &geojson.Feature{
			Geometry: &geojson.Geometry{
				Type:  geojson.GeometryPoint,
				Point: []float64{s.Lng, s.Lat},
			},
			Properties: map[string]interface{}{
				"id":             s.StandID,
				"name":           s.Name,
				"type":           s.Type,
				"numberOfStands": s.NumberOfStands,
				"notes":          s.Notes,
				"imageId":        s.ImageID,
				"source":         s.Source,
				"checked":        s.Checked != "",
				"verified":       s.Verified,
				"publicImageURL": s.PublicImageURL,
			},
		}

		fc.AddFeature(feature)
	}

	err := json.NewEncoder(w).Encode(fc)
	if err != nil {
		fmt.Printf("error encoding json: %s", err)
	}
}

func (a *api) getStand(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("cache-control", "max-age=3600")

	vars := mux.Vars(r)

	s := &stand.Stand{}
	a.DB.Where("`stand_id` = ?", vars["id"]).First(s)

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":             s.StandID,
		"name":           s.Name,
		"type":           s.Type,
		"numberOfStands": s.NumberOfStands,
		"notes":          s.Notes,
		"source":         s.Source,
		"checked":        s.Checked != "",
		"verified":       s.Verified,
		"publicImageURL": s.PublicImageURL,
	})
	if err != nil {
		fmt.Printf("error encoding json: %s", err)
	}
}

func (a *api) createStand(w http.ResponseWriter, r *http.Request) {
	var stand stand.Stand

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			fmt.Printf("error writing error: %s", err)
			return
		}
		return
	}

	err = json.Unmarshal(body, &stand)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("error writing error: %s", err)
		return
	}

	err = validator.Validate(&stand)
	if err != nil {
		errs := err.(validator.ErrorMap)
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errs,
		})
		if err != nil {
			fmt.Printf("error writing error: %s", err)
			return
		}
		return
	}

	stand.Checked = ""
	stand.Source = "User Submission"
	id := uuid.New().String()[:6]
	stand.SourceID = id
	stand.StandID = id
	stand.Token = uuid.New().String()

	a.DB.Create(&stand)

	if a.SendgridAPIKey != "" {
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
	}

	if a.Slack != nil {
		go func() {
			err := a.Slack.PostNotification(stand)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	err = json.NewEncoder(w).Encode(stand)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("could not write json: %s", err)
		return
	}
}

func (a *api) updateStand(w http.ResponseWriter, r *http.Request) {
	var originStand, updatedStand stand.Stand

	vars := mux.Vars(r)

	a.DB.Where("stand_id = ?", vars["id"]).First(&originStand)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("could not write json: %s", err)
			return
		}
		return
	}
	err = json.Unmarshal(body, &updatedStand)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("could not unmarshal json: %s", err)
		return
	}

	err = validator.Validate(&updatedStand)
	if err != nil {
		errs := err.(validator.ErrorMap)
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": errs,
		})
		if err != nil {
			fmt.Printf("error encoding json: %s", err)
			return
		}
		return
	}

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
			"last_update_by":   sub,
		}).Error
		if err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("could not write error: %s", err)
				return
			}
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
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("could not write error: %s", err)
				return
			}
			return
		}
		update := StandUpdate{
			Stand:     originStand,
			UserEmail: sub,
			Update:    string(updateJson),
		}
		if err := tx.Create(&update).Error; err != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("could not write error: %s", err)
				return
			}
			return
		}
		tx.Commit()
	} else {
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []string{
				"stand not found",
			},
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("could not write json: %s", err)
			return
		}
	}

	err = json.NewEncoder(w).Encode(originStand)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("could not write json: %s", err)
		return
	}
}

func (a *api) deleteStand(w http.ResponseWriter, r *http.Request) {
	var stand stand.Stand

	vars := mux.Vars(r)

	a.DB.Where("stand_id = ?", vars["id"]).First(&stand)

	a.DB.Delete(&stand)

	w.WriteHeader(http.StatusAccepted)
}
