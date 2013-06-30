package main

import (
	"github.com/pjvds/httpcallback.io/tests"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

// Setup the test suite
var _ = Suite(&tests.ApiIntegrationTestSuite{
	ProcessFilename: "httpcallback.io",
	ApiBaseUrl:      "http://localhost:8000",
})
