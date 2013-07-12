package tests

import (
	"bytes"
	"fmt"
	"github.com/pjvds/httpcallback.io/api/host"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// The state for the test suite
type ApiIntegrationTestSuite struct {
	ProcessFilename string
	Warmup          int
	ApiBaseUrl      string
	testServer      *httptest.Server
}

// Runs before the test suite starts
func (s *ApiIntegrationTestSuite) SetUpSuite(c *C) {
	configPath := "config.toml"
	config, err := host.OpenConfig(configPath)
	if err != nil {
		c.Fatalf("Cannot open config at %v: %v", configPath, err)
	}

	apiServer := host.NewServer(config)
	testServer := httptest.NewServer(apiServer)
	testServer.Start()
	s.ApiBaseUrl = testServer.URL
	s.testServer = testServer

	c.Logf("Started test server at %v", s.ApiBaseUrl)
}

// Runs after the test suite finished, even when failed
func (s *ApiIntegrationTestSuite) TearDownSuite(c *C) {
	c.Log("Tearing down api integration test suite")

	s.testServer.Close()
	c.Log("Test server closed")
}

func (s *ApiIntegrationTestSuite) TestPing(c *C) {
	response, err := http.Get(s.ApiBaseUrl + "/ping")

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

	response, err := http.Post(s.ApiBaseUrl+"/users", "application/json", buf)

	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusOK)

	doc, err := GetBodyAsDocument(response)

	c.Assert(err, IsNil)
	c.Assert(doc["id"], NotNil)
	c.Assert(len(doc["id"].(string)), Equals, 24)
	c.Assert(doc["username"], Equals, user["username"])
	c.Assert(doc["authToken"], NotNil)
}

func (s *ApiIntegrationTestSuite) TestPostNewUserGetsActuallyAdded(c *C) {
	user := Document{
		"username": "pjvds",
		"password": "foobar",
		"email":    "pj@born2code.net:",
	}
	data := user.ToJson()
	buf := bytes.NewBuffer(data)

	response, err := http.Post(s.ApiBaseUrl+"/users", "application/json", buf)

	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusOK)

	creationReponse, _ := GetBodyAsDocument(response)
	response, err = http.Get(fmt.Sprintf(s.ApiBaseUrl+"/user/%v", creationReponse["id"]))
	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusOK)

	usersResponse, err := GetBodyAsDocument(response)

	c.Assert(err, IsNil)
	c.Assert(usersResponse["id"], Equals, creationReponse["id"])
}

// Users need to be authorized before they can post new callbacks
func (s *ApiIntegrationTestSuite) TestPostNewCallbackUnauthorized(c *C) {
	callback := Document{
		"when": "2006-01-02T15:04:05Z",
		"url":  "http://google.com/",
	}
	data := callback.ToJson()
	buf := bytes.NewBuffer(data)

	response, err := http.Post(s.ApiBaseUrl+"/user/callbacks", "application/json", buf)
	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusUnauthorized)
}

func (s *ApiIntegrationTestSuite) TestPostNewCallbackSuccess(c *C) {
	user := Document{
		"username": "addnewcallbackuser",
		"password": "foobar",
		"email":    "addnewcallbackuser@httpcallback.io",
	}
	data := user.ToJson()
	buf := bytes.NewBuffer(data)

	response, err := http.Post(s.ApiBaseUrl+"/users", "application/json", buf)
	doc, _ := GetBodyAsDocument(response)
	authToken := doc["authToken"]

	callback := Document{
		"when": "2006-01-02T15:04:05Z",
		"url":  "http://google.com/",
	}
	data = callback.ToJson()
	buf = bytes.NewBuffer(data)

	rawUrl := fmt.Sprintf(s.ApiBaseUrl+"/user/callbacks?auth_username=%v&auth_token=%v",
		url.QueryEscape(user["username"].(string)), url.QueryEscape(authToken.(string)))
	response, err = http.Post(rawUrl, "application/json", buf)

	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusOK)
}

func (s *ApiIntegrationTestSuite) TestGetUserReturnsStatusNotFound(c *C) {
	response, err := http.Get(s.ApiBaseUrl + "/user/123")
	c.Assert(err, IsNil)
	c.Assert(response.StatusCode, Equals, http.StatusNotFound)
}
