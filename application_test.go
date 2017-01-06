package ccv2_test

import (
	"context"
	"fmt"
	"net/http"

	. "github.com/Bo0mer/ccv2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Application", func() {
	var client *Client
	var server *ghttp.Server

	var err error

	BeforeEach(func() {
		client, server = setupTestClientAndServer()
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Applications", func() {

		var applications []Application

		JustBeforeEach(func() {
			applications, err = client.Applications(context.Background())
		})

		Context("when the server returns a valid response", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v2/apps"),
						ghttp.RespondWith(http.StatusOK, `
{
    "next_url": null,
    "resources": [
        {
            "entity": {
				"buildpack": "buildpack",
				"command": "command",
				"detected_buildpack": "detected_buildpack",
                "detected_start_command": "detected_start_command",
				"diego": true,
                "disk_quota": 1024,
                "enable_ssh": true,
                "health_check_timeout": 30,
                "health_check_type": "port",
                "instances": 1,
                "memory": 1024,
                "name": "name-2443",
                "package_state": "PENDING",
                "package_updated_at": "2016-06-08T16:41:45Z",
                "space_guid": "9c5c8a91-a728-4608-9f5e-6c8026c3a2ac",
                "stack_guid": "f6c960cc-98ba-4fd1-b197-ecbf39108aa2",
                "staging_failed_description": null,
                "staging_failed_reason": null,
                "staging_task_id": null,
                "state": "STOPPED",
                "version": "f5696e0f-087d-49b0-9ad7-4756c49a6ba6"
            },
            "metadata": {
                "created_at": "2016-06-08T16:41:45Z",
                "guid": "6064d98a-95e6-400b-bc03-be65e6d59622",
                "updated_at": "2016-06-08T16:41:45Z"
            }
        }
    ]
}
					`),
					))
			})

			It("should have returned a list of applications", func() {
				Ω(applications).Should(HaveLen(1))
				app := applications[0]
				Ω(app.GUID).Should(Equal("6064d98a-95e6-400b-bc03-be65e6d59622"))
				Ω(app.CreatedAt).Should(Equal("2016-06-08T16:41:45Z"))
				Ω(app.UpdatedAt).Should(Equal("2016-06-08T16:41:45Z"))
				Ω(app.Entity.Name).Should(Equal("name-2443"))
				Ω(app.Entity.SpaceGUID).Should(Equal("9c5c8a91-a728-4608-9f5e-6c8026c3a2ac"))
				Ω(app.Entity.StackGUID).Should(Equal("f6c960cc-98ba-4fd1-b197-ecbf39108aa2"))
				Ω(app.Entity.Memory).Should(Equal(1024))
				Ω(app.Entity.Instances).Should(Equal(1))
				Ω(app.Entity.DiskQuota).Should(Equal(1024))
				Ω(app.Entity.State).Should(Equal("STOPPED"))
				Ω(app.Entity.Version).Should(Equal("f5696e0f-087d-49b0-9ad7-4756c49a6ba6"))
				Ω(app.Entity.PackageState).Should(Equal("PENDING"))
				Ω(app.Entity.HealthCheckType).Should(Equal("port"))
				Ω(app.Entity.HealthCheckTimeout).Should(Equal(30))
				Ω(app.Entity.Buildpack).Should(Equal("buildpack"))
				Ω(app.Entity.Command).Should(Equal("command"))
				Ω(app.Entity.DetectedBuildpack).Should(Equal("detected_buildpack"))
				Ω(app.Entity.DetectedCommand).Should(Equal("detected_start_command"))
				Ω(app.Entity.Diego).Should(BeTrue())
				Ω(app.Entity.EnableSSH).Should(BeTrue())
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

	Describe("ApplicationSummary", func() {

		var app Application
		var summary ApplicationSummary

		JustBeforeEach(func() {
			summary, err = client.ApplicationSummary(context.Background(), app)
		})

		Context("when the server response is valid", func() {
			BeforeEach(func() {
				app.GUID = "d897c8c-3171-456d-b5d7-3c87feeabbd1"
				path := fmt.Sprintf("/v2/apps/%s/summary", app.GUID)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", path),
						ghttp.RespondWith(http.StatusOK, `
{
	"buildpack": "buildpack",
	"command": "command",
	"detected_buildpack": "detected_buildpack",
    "detected_start_command": "detected_start_command",
	"diego": true,
    "disk_quota": 1024,
    "enable_ssh": true,
    "guid": "cd897c8c-3171-456d-b5d7-3c87feeabbd1",
	"health_check_timeout": 30,
    "health_check_type": "port",
    "instances": 1,
    "memory": 1024,
    "name": "name-79",
    "package_state": "PENDING",
    "package_updated_at": "2016-06-08T16:41:22Z",
    "running_instances": 0,
    "space_guid": "1053174d-eb79-4f16-bf82-9f83a52d6e84",
    "stack_guid": "aff73b55-7767-4928-b0ce-502cca863be0",
    "state": "STOPPED",
    "version": "d457b51a-d7cb-494d-b39e-3171ec75bd60"
}
						`)))
			})

			It("should have not returned an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should have returned summary about the application", func() {
				Ω(summary.GUID).Should(Equal("cd897c8c-3171-456d-b5d7-3c87feeabbd1"))
				Ω(summary.Name).Should(Equal("name-79"))
				Ω(summary.SpaceGUID).Should(Equal("1053174d-eb79-4f16-bf82-9f83a52d6e84"))
				Ω(summary.StackGUID).Should(Equal("aff73b55-7767-4928-b0ce-502cca863be0"))
				Ω(summary.Memory).Should(Equal(1024))
				Ω(summary.Instances).Should(Equal(1))
				Ω(summary.DiskQuota).Should(Equal(1024))
				Ω(summary.State).Should(Equal("STOPPED"))
				Ω(summary.Version).Should(Equal("d457b51a-d7cb-494d-b39e-3171ec75bd60"))
				Ω(summary.PackageState).Should(Equal("PENDING"))
				Ω(summary.HealthCheckType).Should(Equal("port"))
				Ω(summary.HealthCheckTimeout).Should(Equal(30))
				Ω(summary.Buildpack).Should(Equal("buildpack"))
				Ω(summary.Command).Should(Equal("command"))
				Ω(summary.DetectedBuildpack).Should(Equal("detected_buildpack"))
				Ω(summary.DetectedCommand).Should(Equal("detected_start_command"))
				Ω(summary.Diego).Should(BeTrue())
				Ω(summary.EnableSSH).Should(BeTrue())
				Ω(summary.RunningInstances).Should(Equal(0))
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
})
