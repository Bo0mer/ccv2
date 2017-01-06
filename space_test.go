package ccv2_test

import (
	"context"
	"net/http"

	. "github.com/Bo0mer/ccv2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Spaces", func() {
	var client *Client
	var server *ghttp.Server

	var queries []Query
	var spaces []Space
	var err error

	BeforeEach(func() {
		client, server = setupTestClientAndServer()
	})

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		spaces, err = client.Spaces(context.Background(), queries...)
	})

	Context("when the server returns a valid response", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/spaces", "q=name%3Arocket"),
					ghttp.RespondWith(http.StatusOK, `
{
    "next_url": null,
    "resources": [
        {
            "entity": {
                "allow_ssh": true,
                "name": "rocket",
                "organization_guid": "d154425c-dccc-42e6-b6b4-27d46c3b42cb",
                "space_quota_definition_guid": null
            },
            "metadata": {
                "created_at": "2016-06-08T16:41:40Z",
                "guid": "2e100106-0b74-4062-8671-0d375f951cb4",
                "updated_at": "2016-06-08T16:41:26Z"
            }
        }
    ]
}`),
				),
			)

			nameQuery := Query{
				Filter: FilterName,
				Op:     OperatorEqual,
				Value:  "rocket",
			}
			queries = []Query{nameQuery}
		})

		It("should have returned a list of spaces", func() {
			Ω(spaces).Should(HaveLen(1))
			space := spaces[0]
			Ω(space.GUID).Should(Equal("2e100106-0b74-4062-8671-0d375f951cb4"))
			Ω(space.CreatedAt).Should(Equal("2016-06-08T16:41:40Z"))
			Ω(space.UpdatedAt).Should(Equal("2016-06-08T16:41:26Z"))
			Ω(space.Entity.Name).Should(Equal("rocket"))
			Ω(space.Entity.OrganizationGUID).Should(Equal("d154425c-dccc-42e6-b6b4-27d46c3b42cb"))
			Ω(space.Entity.SpaceQuotaDefinitonGUID).Should(Equal(""))
			Ω(space.Entity.AllowSSH).Should(Equal(true))
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

})
