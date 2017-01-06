package ccv2_test

import (
	"net/http"
	"net/url"

	"github.com/Bo0mer/ccv2"
	"github.com/onsi/gomega/ghttp"
)

func setupTestClientAndServer() (*ccv2.Client, *ghttp.Server) {
	server := ghttp.NewServer()
	u, err := url.Parse("http://" + server.Addr())
	if err != nil {
		panic(err)
	}
	client := &ccv2.Client{
		API:        u,
		HTTPClient: http.DefaultClient,
	}

	return client, server
}

func notFoundHandler() http.HandlerFunc {
	return ghttp.RespondWith(http.StatusNotFound, `{"error_code":"10001","description":"Entity is missing."}`)
}

var notFoundErr = &ccv2.UnexpectedResponseError{
	StatusCode:  http.StatusNotFound,
	ErrorCode:   "10001",
	Description: "Entity is missing.",
}
