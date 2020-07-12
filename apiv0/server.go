package apiv0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"code.katiechapman.ie/dublinbikeparking/stand"

	"code.katiechapman.ie/dublinbikeparking/slack"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	minio "github.com/minio/minio-go/v6"
	"github.com/osstotalsoft/oidc-jwt-go"
	"github.com/osstotalsoft/oidc-jwt-go/discovery"
)

type api struct {
	DB             *gorm.DB
	SendgridAPIKey string
	Slack          *slack.SlackIntegration
	validator      func(request *http.Request) (*jwt.Token, error)
}

func NewAPIv0(r *mux.Router, db *gorm.DB) {
	apiHandler := &api{
		DB: db,
	}

	if os.Getenv("SENDGRID_API_KEY") == "" {
		fmt.Println("WARNING: SENDGRID_API_KEY is not set. No mails will be sent")
	}
	apiHandler.SendgridAPIKey = os.Getenv("SENDGRID_API_KEY")

	if os.Getenv("S3_ENDPOINT") == "" {
		fmt.Println("WARNING: S3_* variables ares not set. No images will be stored")
	}

	if webhookURL := os.Getenv("SLACK_WEBHOOK_URL"); webhookURL == "" {
		fmt.Println("WARNING: S3_* variables ares not set. No images will be stored")
	} else {
		apiHandler.Slack = slack.NewSlackIntegration(webhookURL)
	}

	db.AutoMigrate(&stand.Stand{})
	db.AutoMigrate(&StandUpdate{})

	var minioClient *minio.Client
	var endpoint, accessKeyID, secretAccessKey, bucketName string
	var useSSL bool

	if os.Getenv("OIDC_AUTHORITY") == "" {
		fmt.Println("ERROR: OIDC_AUTHORITY variable is not set.")
		os.Exit(1)
	}
	if os.Getenv("OIDC_AUDIENCE") == "" {
		fmt.Println("ERROR: OIDC_AUDIENCE variable is not set.")
		os.Exit(1)
	}

	authority := os.Getenv("OIDC_AUTHORITY")
	audience := os.Getenv("OIDC_AUDIENCE")

	secretProvider := oidc.NewOidcSecretProvider(
		discovery.NewClient(discovery.Options{
			Authority: authority,
		}),
	)
	validator := oidc.NewJWTValidator(request.AuthorizationHeaderExtractor, secretProvider, audience, authority)
	apiHandler.validator = validator

	if os.Getenv("S3_ENDPOINT") != "" {
		endpoint = os.Getenv("S3_ENDPOINT")
		accessKeyID = os.Getenv("S3_ACCESS_KEY_ID")
		secretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
		bucketName = os.Getenv("S3_BUCKET_NAME")
		useSSL = true
		if os.Getenv("S3_USE_SSL") == "false" {
			useSSL = false
		}
		var err error
		minioClient, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
		if err != nil {
			log.Fatalln(err)
		}
	}

	r.HandleFunc("/stand", apiHandler.getStands).Methods("GET")
	r.HandleFunc("/stand", apiHandler.createStand).Methods("POST")
	r.HandleFunc("/stand/{id}", apiHandler.getStand).Methods("GET")
	r.Handle("/stand/{id}", apiHandler.authMiddleware(http.HandlerFunc(apiHandler.updateStand))).Methods("POST")
	r.Handle("/stand/{id}", apiHandler.authMiddleware(http.HandlerFunc(apiHandler.deleteStand))).Methods("DELETE")
	r.HandleFunc("/image", handleImageOptionsFunc(minioClient)).Methods("OPTIONS")
	r.HandleFunc("/image", handleImagePostFunc(minioClient, bucketName)).Methods("POST")
	r.HandleFunc("/publicimage/{id}", apiHandler.handlePublicImagePostFunc(minioClient, "dublinbikeparking-public")).Methods("POST")
	r.Handle("/image/{id}", apiHandler.authMiddleware(http.HandlerFunc(handleImageGetFunc(minioClient, bucketName)))).Methods("GET")
	r.HandleFunc("/slack", apiHandler.handleSlackMessage).Methods("POST")
}

func handleImageOptionsFunc(minioClient *minio.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if minioClient == nil {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
func handleImageGetFunc(minioClient *minio.Client, bucketName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if minioClient == nil {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

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
		if minioClient == nil {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

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
		fileExpiration := time.Now().Add(time.Hour * 24 * 14)
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

func (a api) handlePublicImagePostFunc(minioClient *minio.Client, bucketName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if minioClient == nil {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		vars := mux.Vars(r)

		var stnd stand.Stand
		a.DB.Where("stand_id = ?", vars["id"]).First(&stnd)

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

		err = a.DB.Model(&stnd).Limit(1).Update(map[string]interface{}{
			"public_image_url": "https://f003.backblazeb2.com/file/dublinbikeparking-public/" + fileName,
		}).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			fmt.Fprintf(w, "error uploading file")
			return
		}

		fmt.Fprintf(w, "%s", fileName)
	}
}

type SlackInteraction struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	ResponseURL string `json:"response_url"`
	Actions     []struct {
		ActionID string `json:"action_id"`
		Value    string `json:"value"`
	} `json:"actions"`
}

func (a *api) handleSlackMessage(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	interaction := &SlackInteraction{}
	err := json.Unmarshal([]byte(payload), interaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		fmt.Fprintf(w, "error approving slack message")
		return
	}

	standID := (&url.URL{
		RawQuery: interaction.Actions[0].Value,
	}).Query().Get("id")
	standToken := (&url.URL{
		RawQuery: interaction.Actions[0].Value,
	}).Query().Get("token")

	go func() {
		br := &bytes.Buffer{}
		err = json.NewEncoder(br).Encode(struct {
			ReplaceOriginal string `json:"replace_original"`
			Text            string `json:"text"`
		}{
			ReplaceOriginal: "true",
			Text:            fmt.Sprintf("<@%s> has %s stand ID %s", interaction.User.ID, interaction.Actions[0].ActionID, standID),
		})
		if err != nil {
			log.Println(err)
			return
		}
		_, err = http.Post(interaction.ResponseURL, "application/json", br)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	s := stand.Stand{}

	var errSw error
	switch interaction.Actions[0].ActionID {
	case "approve":
		errSw = a.DB.Model(&s).Where("stand_id = ? AND token = ?", standID, standToken).Update("checked", interaction.User.ID).Error
	default:
		errSw = a.DB.Model(&s).Where("stand_id = ? AND token = ?", standID, standToken).Update("deleted_at", time.Now()).Error
	}
	if errSw != nil {
		log.Println(errSw)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := a.validator(r)
		if err != nil {
			log.Println("AuthorizationFilter: Token is not valid", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		context.Set(r, "userSub", claims["sub"])

		next.ServeHTTP(w, r)
	})
}
