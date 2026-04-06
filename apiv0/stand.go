package apiv0

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/honeycombio/beeline-go"

	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	geojson "github.com/paulmach/go.geojson"
	"gopkg.in/validator.v2"

	"github.com/semanticallynull/dublinbikeparking/stand"
)

type StandUpdate struct {
	gorm.Model
	Stand     stand.Stand
	StandID   uint
	UserEmail string
	Update    string `sql:"type:text"`
}

type geoJSONCache struct {
	featureCollection *geojson.FeatureCollection
	expiry            time.Time
}

var cache = geoJSONCache{}

func (a *api) getStands(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Add("content-type", "application/json")
	fc := geojson.NewFeatureCollection()

	checked := r.URL.Query().Get("checked") != "unchecked"
	db := r.URL.Query().Get("dublinbikes") == "only"
	review := r.URL.Query().Get("review") == "true"

	if checked && !db && !review && cache.expiry.After(time.Now()) {
		beeline.AddField(ctx, "cached", true)
		err := json.NewEncoder(w).Encode(cache.featureCollection)
		if err != nil {
			fmt.Printf("error encoding json: %s", err)
		}
		return
	}

	ctx, span := beeline.StartSpan(ctx, "getStands: db")
	beeline.AddField(ctx, "cached", false)

	dbc := a.DB

	if checked {
		dbc = a.DB.Where("checked != ?", "")
	}

	if db {
		dbc = dbc.Where("type = ?", "DublinBikes")
	} else {
		dbc = dbc.Where("type != ?", "DublinBikes")
	}

	if review {
		dbc = dbc.Where("number_of_stands IS NULL")
	}

	stands := []stand.Stand{}
	dbc.Find(&stands)

	span.Send()

	ctx, span = beeline.StartSpan(ctx, "getStands: toGeoJSON")

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

	if checked && !db && !review {
		cache.featureCollection = fc
		cache.expiry = time.Now().Add(time.Hour)
	}

	span.Send()

	err := json.NewEncoder(w).Encode(fc)
	if err != nil {
		fmt.Printf("error encoding json: %s", err)
	}
}

func (a *api) standMissing(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	vars := mux.Vars(r)

	// If POST, check password for authenticated missing reports (verify mode)
	if r.Method == "POST" {
		var body struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}
		expected := os.Getenv("DBP_VERIFY_PASSWORD")
		if expected == "" || body.Password != expected {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid password"})
			return
		}
	}

	s := stand.Stand{}
	if query := a.DB.Where("`stand_id` = ?", vars["id"]).First(&s); query.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(query.Error)
		return
	}

	if s.Token == "" {
		token := uuid.New().String()
		if query := a.DB.Model(&s).Where("`stand_id` = ?", vars["id"]).Update("token", token); query.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(query.Error)
			return
		}
	}

	if a.SMTP != nil {
		go func() {
			subject := fmt.Sprintf("Stand Missing at '%s'", s.Name)
			body := "Stand missing DublinBikeParking.com\n" +
				"Stand ID: " + s.StandID + "\n" +
				"Name: " + s.Name + "\n" +
				"Coordinates: " + fmt.Sprintf("%f %f", s.Lat, s.Lng) + "\n" +
				"https://www.dublinbikeparking.com/#19/" + fmt.Sprintf("%f/%f", s.Lat, s.Lng)
			if err := sendMail(a.SMTP, subject, body); err != nil {
				log.Println("email error:", err)
			}
		}()
	}

	if a.Slack != nil {
		go func() {
			err := a.Slack.PostMissingNotification(s)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	cache = geoJSONCache{}

	err := json.NewEncoder(w).Encode("OK")
	if err != nil {
		fmt.Printf("error encoding json: %s", err)
	}
}

func (a *api) standVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	vars := mux.Vars(r)

	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	expected := os.Getenv("DBP_VERIFY_PASSWORD")
	if expected == "" || body.Password != expected {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid password"})
		return
	}

	s := stand.Stand{}
	if query := a.DB.Where("`stand_id` = ?", vars["id"]).First(&s); query.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(query.Error)
		return
	}

	if err := a.DB.Model(&s).Where("`stand_id` = ?", vars["id"]).Updates(map[string]interface{}{
		"checked":  "rider-verify",
		"verified": true,
	}).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	cache = geoJSONCache{}

	if a.Slack != nil {
		go func() {
			err := a.Slack.PostVerifyNotification(s)
			if err != nil {
				log.Println("slack verify notification error:", err)
			}
		}()
	}

	err := json.NewEncoder(w).Encode("OK")
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

	standData := map[string]interface{}{
		"id":             s.StandID,
		"name":           s.Name,
		"type":           s.Type,
		"numberOfStands": s.NumberOfStands,
		"notes":          s.Notes,
		"source":         s.Source,
		"checked":        s.Checked != "",
		"verified":       s.Verified,
		"publicImageURL": s.PublicImageURL,
	}

	if s.Type == "DublinBikes" {
		dbAPIKey := os.Getenv("DUBLINBIKES_API_KEY")
		if dbAPIKey != "" {
			resp, err := http.Get(fmt.Sprintf("https://api.jcdecaux.com/vls/v1/stations/%s?contract=dublin&apiKey=%s", s.SourceID, dbAPIKey))
			if err == nil && resp.StatusCode == 200 {
				ld := LiveData{}
				json.NewDecoder(resp.Body).Decode(&ld)
				standData["bikesAvailable"] = ld.AvailableBikes
				standData["standsAvailable"] = ld.AvailableBikeStands
			}
		}
	}

	err := json.NewEncoder(w).Encode(standData)
	if err != nil {
		fmt.Printf("error encoding json: %s", err)
	}
}

type LiveData struct {
	AvailableBikes      int `json:"available_bikes"`
	AvailableBikeStands int `json:"available_bike_stands"`
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

	if a.SMTP != nil {
		go func() {
			subject := fmt.Sprintf("New Stand at '%s'", stand.Name)
			body := "New stand on DublinBikeParking.com\n" +
				"Stand ID: " + stand.StandID + "\n" +
				"Name: " + stand.Name + "\n" +
				"Coordinates: " + fmt.Sprintf("%f %f", stand.Lat, stand.Lng) + "\n" +
				"https://www.dublinbikeparking.com/#19/" + fmt.Sprintf("%f/%f", stand.Lat, stand.Lng)
			if err := sendMail(a.SMTP, subject, body); err != nil {
				log.Println("email error:", err)
			}
		}()
	}

	if a.Slack != nil {
		go func() {
			err := a.Slack.PostCreateNotification(stand)
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

	body, err := io.ReadAll(r.Body)
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
