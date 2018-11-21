package integration_test

import (
	"net"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	binPath    string
	serverPort string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	var err error
	binPath, err = gexec.Build("code.benchapman.ie/dublinbikeparking")
	Expect(err).ToNot(HaveOccurred())

	os.Setenv("PORT", getFreePort())

	execBin()
	time.Sleep(time.Second)
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func execBin(args ...string) *gexec.Session {
	cmd := exec.Command(binPath, args...)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	return session
}

func getFreePort() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	serverPort = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	listener.Close()
	return serverPort
}
