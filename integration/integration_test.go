package integration_test

import (
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
)

var _ = Describe("DublinBikeParking", func() {
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
