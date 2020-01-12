package apiv0

import (
	"fmt"
	"net/http"
	"os"

	"github.com/auth0-community/go-auth0"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gopkg.in/square/go-jose.v2"
)

type api struct {
	DB             *gorm.DB
	SendgridAPIKey string
}

func NewAPIv0(r *mux.Router, db *gorm.DB) {
	apiHandler := &api{
		DB: db,
	}

	if os.Getenv("SENDGRID_API_KEY") == "" {
		fmt.Println("You must set a SENDGRID_API_KEY")
		os.Exit(1)
	}

	apiHandler.SendgridAPIKey = os.Getenv("SENDGRID_API_KEY")

	db.AutoMigrate(&Stand{})
	db.AutoMigrate(&StandUpdate{})
	db.AutoMigrate(&Theft{})

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
		context.Set(r, "userSub", out["sub"])

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
