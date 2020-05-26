package apiv0

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	minio "github.com/minio/minio-go/v6"
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

	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	useSSL := true
	if os.Getenv("S3_USE_SSL") == "false" {
		useSSL = false
	}

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	r.HandleFunc("/stand", apiHandler.getStands).Methods("GET")
	r.HandleFunc("/stand", apiHandler.createStand).Methods("POST")
	r.HandleFunc("/stand/{id}", apiHandler.getStand).Methods("GET")
	r.Handle("/stand/{id}", authMiddleware(http.HandlerFunc(apiHandler.updateStand))).Methods("POST")
	r.Handle("/stand/{id}", authMiddleware(http.HandlerFunc(apiHandler.deleteStand))).Methods("DELETE")
	r.HandleFunc("/image", handleImagePostFunc(minioClient, bucketName)).Methods("POST")
	r.Handle("/image/{id}", authMiddleware(http.HandlerFunc(handleImageGetFunc(minioClient, bucketName)))).Methods("GET")
}

func handleImageGetFunc(minioClient *minio.Client, bucketName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		_, err := minioClient.StatObject(bucketName, vars["id"], minio.StatObjectOptions{})
		if err != nil {
			switch minio.ToErrorResponse(err).StatusCode {
			case 404:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Println(err)
			return
		}

		reqParams := make(url.Values)
		presignedURL, err := minioClient.PresignedGetObject(bucketName, vars["id"], time.Minute*15, reqParams)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", presignedURL)
	}
}

func handleImagePostFunc(minioClient *minio.Client, bucketName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(5 << 20)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			fmt.Fprintf(w, "Image too large. Max Size: %v", 5<<20)
			return
		}

		file, fileHeader, err := r.FormFile("filepond")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			fmt.Fprintf(w, "error uploading file")
			return
		}
		defer file.Close()

		contentType := fileHeader.Header.Get("Content-Type")
		if !(contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif") {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "incorrect filetype")
			return
		}

		fileName := uuid.New().String()
		_, err = minioClient.PutObject(bucketName, fileName, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			fmt.Fprintf(w, "error uploading file")
			return
		}

		fmt.Fprintf(w, "%s", fileName)
	}
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
