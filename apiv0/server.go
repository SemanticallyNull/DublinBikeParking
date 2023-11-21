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

	"github.com/gin-gonic/gin"

	"github.com/semanticallynull/dublinbikeparking/slack"
	"github.com/semanticallynull/dublinbikeparking/stand"

	"github.com/honeycombio/beeline-go"
	"github.com/jinzhu/gorm"
	"github.com/minio/minio-go/v6"
)

type api struct {
	DB             *gorm.DB
	SendgridAPIKey string
	Slack          *slack.SlackIntegration
}

func NewAPIv0(r *gin.RouterGroup, db *gorm.DB) {
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

	if honeycombWriteKey := os.Getenv("DBP_HONEYCOMB_WRITEKEY"); honeycombWriteKey != "" {
		beeline.Init(beeline.Config{
			WriteKey:    honeycombWriteKey,
			Dataset:     "DublinBikeParking",
			ServiceName: "dbp",
		})
	}

	db.AutoMigrate(&stand.Stand{})
	db.AutoMigrate(&StandUpdate{})

	var minioClient *minio.Client
	var endpoint, accessKeyID, secretAccessKey, bucketName string
	var useSSL bool

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

	r.GET("/hirebikes", wrap(apiHandler.getHireBikes))
	r.GET("/stand", wrap(apiHandler.getStands))
	r.POST("/stand", wrap(apiHandler.createStand))
	r.GET("/stand/{id}", wrap(apiHandler.getStand))
	r.GET("/stand/{id}/missing", wrap(apiHandler.standMissing))
	r.OPTIONS("/image", wrap(handleImageOptionsFunc(minioClient)))
	r.POST("/image", wrap(handleImagePostFunc(minioClient, bucketName)))
	r.POST("/publicimage/{id}", wrap(apiHandler.handlePublicImagePostFunc(minioClient, "dublinbikeparking-public")))
	r.POST("/slack", wrap(apiHandler.handleSlackMessage))
}

func wrap(handlerFunc http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handlerFunc(c.Writer, c.Request)
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
	case "dont_hide":
		break
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
