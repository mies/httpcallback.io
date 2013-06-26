package tests

import (
	"bytes"
	. "launchpad.net/gocheck"
	"net/http"
	"os"
	"time"
)

// The state for the test suite
type ApiIntegrationTestSuite struct {
	ProcessFilename string
	process         *os.Process
}

// Runs before the test suite starts
func (s *ApiIntegrationTestSuite) SetUpSuite(c *C) {
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	process, err := os.StartProcess(s.ProcessFilename, []string{"-config config.toml -port 8000"}, &procAttr)
	if err != nil {
		c.Errorf("Unable to start %s: %s", s.ProcessFilename, err.Error())
		c.Fail()
	}

	// Allow process to warm up
	time.Sleep(750 * time.Millisecond)
	s.process = process

	c.Logf("Started %s, pid %v", s.ProcessFilename, process.Pid)
}

// Runs after the test suite finished, even when failed
func (s *ApiIntegrationTestSuite) TearDownSuite(c *C) {
	if err := s.process.Kill(); err != nil {
		c.Logf("Unable to kill %s: %s", s.ProcessFilename, err.Error())
	}
}

func (s *ApiIntegrationTestSuite) TestPing(c *C) {
	response, err := http.Get("http://api.localhost:8000/ping")

	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusOK)

	doc, err := GetBodyAsDocument(response)

	c.Assert(err, IsNil)
	c.Assert(doc["message"], Equals, "pong")
}

func (s *ApiIntegrationTestSuite) TestPostNewUserResponse(c *C) {
	user := Document{
		"username": "pjvds",
		"password": "foobar",
		"email":    "pj@born2code.net:",
	}
	data := user.ToJson()
	buf := bytes.NewBuffer(data)

	response, err := http.Post("http://api.localhost:8000/users", "application/json", buf)

	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusOK)

	doc, err := GetBodyAsDocument(response)

	c.Assert(err, IsNil)
	c.Assert(doc["username"], Equals, user["username"])
}

func (s *ApiIntegrationTestSuite) TestPostNewUserGetsActuallyAdded(c *C) {
	user := Document{
		"username": "pjvds",
		"password": "foobar",
		"email":    "pj@born2code.net:",
	}
	data := user.ToJson()
	buf := bytes.NewBuffer(data)

	response, err := http.Post("http://api.localhost:8000/users", "application/json", buf)

	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusOK)

	creationReponse, err := GetBodyAsDocument(response)

	response, err = http.Get("http://api.localhost:8000/users")
	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusOK)

	usersResponse, err := GetBodyAsDocument(response)

	c.Assert(err, IsNil)
	c.Assert(usersResponse["users"].(map[string]interface{})["Id"], Equals, creationReponse["Id"])
}

func (s *ApiIntegrationTestSuite) TestGetUserReturnsStatusNotFound(c *C) {
	response, err := http.Get("http://api.localhost:8000/user/123")
	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusNotFound)
	c.Assert(response.Status, Equals, "user not found")
}
