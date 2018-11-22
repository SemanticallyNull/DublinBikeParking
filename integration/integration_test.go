package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	geojson "github.com/paulmach/go.geojson"

	"net/http"
)

var _ = Describe("DublinBikeParking", func() {
	Describe("frontend", func() {
		It("returns a 200 for the root", func() {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/", serverPort))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("returns a 200 for the stand_icon_sheffield.png", func() {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/stand_icon_sheffield.png", serverPort))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("returns a 200 for the BleeperData.json", func() {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/BleeperData.json", serverPort))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("returns a 200 for the healthz endpoint", func() {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/healthz", serverPort))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte("ok")))
		})
	})

	Describe("api", func() {
		It("returns json when I get stands", func() {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v0/stand", serverPort))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			output := geojson.FeatureCollection{}
			decodeErr := json.NewDecoder(resp.Body).Decode(&output)
			Expect(decodeErr).ToNot(HaveOccurred())
		})
	})
})
