# ccv2
[![GoDoc](https://godoc.org/github.com/Bo0mer/ccv2?status.svg)](https://godoc.org/github.com/Bo0mer/ccv2)
[![Build Status](https://travis-ci.org/Bo0mer/ccv2.svg?branch=master)](https://travis-ci.org/Bo0mer/ccv2)


Package ccv2 provides read-only Cloud Foundry Cloud Controller API client.
The client targets version 2 of the Cloud Controller's API.

For more details on the API itself, please refer to https://apidocs.cloudfoundry.org.
For more details on the client's Go API, refer to its [godoc](https://godoc.org/github.com/Bo0mer/ccv2).

## Installation
As most Go packages, just go get it.
```
go get -u github.com/Bo0mer/ccv2
```

## Developer's guide
If you want to introduce new functionality, fix a bug, or just play with the code,
just make sure the tests are passing by executing `go test -race ./...` or `ginkgo -r --race`.
