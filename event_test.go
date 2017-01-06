package ccv2_test

import (
	"context"
	"net/http"
	"time"

	. "github.com/Bo0mer/ccv2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Events", func() {
	var client *Client
	var server *ghttp.Server

	var queries []Query
	var events []Event
	var err error

	BeforeEach(func() {
		client, server = setupTestClientAndServer()
	})

	JustBeforeEach(func() {
		events, err = client.Events(context.Background(), queries...)
	})

	Context("when the server returns a valid response", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/events", "q=type%3Aapp.crash"),
					ghttp.RespondWith(http.StatusOK, `
{
    "next_url": null,
    "resources": [
        {
            "metadata": {
                "created_at": "2016-06-08T16:41:23Z",
                "guid": "b8ede8e1-afc8-40a1-baae-236a0a77b27b",
                "updated_at": "2016-06-08T16:41:26Z"
            },
            "entity": {
                "actee": "guid-e7790fa4-be2b-4a0f-aa82-c124342b0bb4",
                "actee_name": "name-171",
                "actee_type": "name-170",
                "actor": "guid-008640fc-d316-4602-9251-c8d09bbdc750",
                "actor_name": "name-169",
                "actor_type": "name-168",
                "organization_guid": "86aa12ee-8c4f-4b26-b391-2be6c1730dbc",
                "space_guid": "3a1368e7-e3b7-46af-a98d-57b9c71445e7",
                "timestamp": "2016-06-08T16:41:23Z",
                "type": "app.crash"
            }
        }
    ]
}`),
				),
			)

			typeQuery := Query{
				Filter: FilterType,
				Op:     OperatorEqual,
				Value:  "app.crash",
			}
			queries = []Query{typeQuery}
		})

		It("should have not returned an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should have returned list of events", func() {
			Ω(events).Should(HaveLen(1))
			event := events[0]
			Ω(event.GUID).Should(Equal("b8ede8e1-afc8-40a1-baae-236a0a77b27b"))
			Ω(event.CreatedAt).Should(Equal("2016-06-08T16:41:23Z"))
			Ω(event.UpdatedAt).Should(Equal("2016-06-08T16:41:26Z"))
			Ω(event.Entity.Actee).Should(Equal("guid-e7790fa4-be2b-4a0f-aa82-c124342b0bb4"))
			Ω(event.Entity.ActeeName).Should(Equal("name-171"))
			Ω(event.Entity.ActeeType).Should(Equal("name-170"))
			Ω(event.Entity.Actor).Should(Equal("guid-008640fc-d316-4602-9251-c8d09bbdc750"))
			Ω(event.Entity.ActorName).Should(Equal("name-169"))
			Ω(event.Entity.ActorType).Should(Equal("name-168"))
			Ω(event.Entity.Type).Should(Equal("app.crash"))
			Ω(event.Entity.SpaceGUID).Should(Equal("3a1368e7-e3b7-46af-a98d-57b9c71445e7"))
			Ω(event.Entity.OrganizationGUID).Should(Equal("86aa12ee-8c4f-4b26-b391-2be6c1730dbc"))

			var timestamp time.Time
			perr := timestamp.UnmarshalText([]byte("2016-06-08T16:41:23Z"))
			Ω(perr).ShouldNot(HaveOccurred())
			Ω(event.Entity.Timestamp).Should(Equal(timestamp))
		})
	})

})
