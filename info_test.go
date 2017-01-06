package ccv2_test

import (
	"context"
	"net/http"

	. "github.com/Bo0mer/ccv2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Info", func() {
	var client *Client
	var server *ghttp.Server

	var info Info
	var err error

	BeforeEach(func() {
		client, server = setupTestClientAndServer()
	})

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		info, err = client.Info(context.Background())
	})

	Context("when the server returns a valid response", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/info"),
					ghttp.RespondWith(http.StatusOK, `
{
   "name": "CF",
   "build": "dev",
   "support": "https://help.bosh-lite.com",
   "version": 249,
   "description": "Cloud Foundry",
   "authorization_endpoint": "https://login.bosh-lite.com",
   "token_endpoint": "https://uaa.bosh-lite.com",
   "min_cli_version": null,
   "api_version": "2.65.0",
   "app_ssh_endpoint": "ssh.bosh-lite.com:2222",
   "logging_endpoint": "wss://loggregator.bosh-lite.com:443",
   "doppler_logging_endpoint": "wss://doppler.bosh-lite.com:443"
}
					`),
				))
		})

		It("should have not returned an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should have returned info", func() {
			Ω(info.Name).Should(Equal("CF"))
			Ω(info.Build).Should(Equal("dev"))
			Ω(info.Support).Should(Equal("https://help.bosh-lite.com"))
			Ω(info.Version).Should(Equal(249))
			Ω(info.Description).Should(Equal("Cloud Foundry"))
			Ω(info.AuthorizationEndpoint).Should(Equal("https://login.bosh-lite.com"))
			Ω(info.TokenEndpoint).Should(Equal("https://uaa.bosh-lite.com"))
			Ω(info.MinCliVersion).Should(Equal(""))
			Ω(info.APIVersion).Should(Equal("2.65.0"))
			Ω(info.AppSSHEndpoint).Should(Equal("ssh.bosh-lite.com:2222"))
			Ω(info.LoggingEndpoint).Should(Equal("wss://loggregator.bosh-lite.com:443"))
			Ω(info.DopplerEndpoint).Should(Equal("wss://doppler.bosh-lite.com:443"))
		})
	})

	Context("when the server returns a non-2XX response", func() {
		BeforeEach(func() {
			server.AppendHandlers(notFoundHandler())
		})

		It("should have returned a UnexpectedResponseError", func() {
			Ω(err).Should(Equal(notFoundErr))
		})
	})

	Context("when the server returns invalid response", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, "-"))
		})

		It("should have returned an error", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

})
