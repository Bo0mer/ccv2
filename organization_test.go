package ccv2_test

import (
	"context"
	"net/http"

	. "github.com/Bo0mer/ccv2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Organizations", func() {
	var client *Client
	var server *ghttp.Server

	var queries []Query
	var organizations []Organization
	var err error

	BeforeEach(func() {
		client, server = setupTestClientAndServer()
		queries = nil
	})

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		organizations, err = client.Organizations(context.Background(), queries...)
	})

	Context("when the server returns a valid response", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/organizations", "q=name%3ANASA"),
					ghttp.RespondWith(http.StatusOK, `
{
    "next_url": "",
    "resources": [
        {
            "entity": {
                "billing_enabled": true,
                "name": "NASA",
                "quota_definition_guid": "dcb680a9-b190-4838-a3d2-b84aa17517a6",
                "status": "active"
            },
            "metadata": {
                "created_at": "2016-06-08T16:41:33Z",
                "guid": "a7aff246-5f5b-4cf8-87d8-f316053e4a20",
                "updated_at": "2016-06-08T16:41:37Z"
            }
        }
    ]
}`),
				),
			)

			nameQuery := Query{
				Filter: FilterName,
				Op:     OperatorEqual,
				Value:  "NASA",
			}
			queries = []Query{nameQuery}
		})

		It("should  have returned list of organizations", func() {
			Ω(organizations).Should(HaveLen(1))
			org := organizations[0]
			Ω(org.GUID).Should(Equal("a7aff246-5f5b-4cf8-87d8-f316053e4a20"))
			Ω(org.CreatedAt).Should(Equal("2016-06-08T16:41:33Z"))
			Ω(org.UpdatedAt).Should(Equal("2016-06-08T16:41:37Z"))
			Ω(org.Entity.Name).Should(Equal("NASA"))
			Ω(org.Entity.BillingEnabled).Should(BeTrue())
			Ω(org.Entity.Status).Should(Equal("active"))
		})

		It("should have not returned an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
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

	Context("when the response is paginated", func() {

		BeforeEach(func() {
			var org1, org2 Organization
			org1.Metadata.GUID = "org-1"
			org2.Metadata.GUID = "org-2"
			p1 := orgPaginatedResponse{
				NextURL:   "http://" + server.Addr() + "/v2/organizations?page=2",
				Resources: []Organization{org1},
			}
			p2 := orgPaginatedResponse{
				NextURL:   "",
				Resources: []Organization{org2},
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/organizations"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, p1)),

				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/organizations", "page=2"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, p2)))
		})

		It("should have returned all spaces", func() {
			Ω(organizations).Should(HaveLen(2))
			Ω(organizations[0].GUID).Should(Equal("org-1"))
			Ω(organizations[1].GUID).Should(Equal("org-2"))
		})
	})
})

type orgPaginatedResponse struct {
	NextURL   string         `json:"next_url"`
	Resources []Organization `json:"resources"`
}
