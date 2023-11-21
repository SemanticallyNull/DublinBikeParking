package apiv0

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v6"

	"github.com/semanticallynull/dublinbikeparking/stand"
)

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
