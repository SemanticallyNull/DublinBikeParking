package apiv0

import (
	"fmt"
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	jose "gopkg.in/square/go-jose.v2"
)

type api struct {
	DB *gorm.DB
}

func NewAPIv0(r *mux.Router, db *gorm.DB) {
	apiHandler := &api{
		DB: db,
	}

	db.AutoMigrate(&Stand{})
	db.AutoMigrate(&StandUpdate{})

	r.HandleFunc("/stand", apiHandler.getStands).Methods("GET")
	r.HandleFunc("/stand", apiHandler.createStand).Methods("POST")
	r.Handle("/stand/{id}", authMiddleware(http.HandlerFunc(apiHandler.updateStand))).Methods("POST")
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		audience := []string{"0Hz3wMPMskh2qVpypXPzjwOykkYV1eZd"}
		secretProvider := auth0.NewJWKClient(auth0.JWKClientOptions{URI: "https://benchapman.eu.auth0.com/.well-known/jwks.json"}, nil)

		configuration := auth0.NewConfiguration(secretProvider, audience, "https://benchapman.eu.auth0.com/", jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)
		out := map[string]interface{}{}
		err = validator.Claims(r, token, &out)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
		context.Set(r, "userEmail", out["email"])

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
