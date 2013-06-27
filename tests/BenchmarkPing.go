package tests

import (
	. "launchpad.net/gocheck"
	"net/http"
)

func (s *ApiIntegrationTestSuite) BenchmarkPing(c *C) {
	for i := 0; i < c.N; i++ {
		response, err := http.Get("http://api.localhost:8000/ping")

		if err != nil {
			c.Fatalf("Error while getting response: %v", err.Error())
		} else {
			response.Body.Close()
		}
	}
}
