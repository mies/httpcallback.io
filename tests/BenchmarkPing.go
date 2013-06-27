package tests

import (
	"bytes"
	"fmt"
	. "launchpad.net/gocheck"
	"net/http"
	"time"
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

func (s *ApiIntegrationTestSuite) BenchmarkNewUser(c *C) {
	t := time.Now().UnixNano()

	for i := 0; i < c.N; i++ {
		user := Document{
			"username": fmt.Sprintf("pjvds%v", int64(i)+t),
			"password": fmt.Sprintf("foobar%v", i),
			"email":    fmt.Sprintf("pj%v@born2code.net", i),
		}
		data := user.ToJson()
		buf := bytes.NewBuffer(data)

		response, err := http.Post("http://api.localhost:8000/users", "application/json", buf)

		if err != nil {
			c.Fatalf("Error while getting response: %v", err.Error())
		} else {
			response.Body.Close()
		}
	}
}
