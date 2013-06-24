package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
	"os"
	"testing"
	"time"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type Document map[string]interface{}

// Setup the test suite
var _ = Suite(&ApiIntegrationTestSuite{
	ProcessFilename: "httpcallback.io",
})

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
	time.Sleep(2 * time.Second)
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

	doc, err := GetBodyAsDocument(c, response)

	c.Assert(err, IsNil)
	c.Assert(doc["message"], Equals, "pong")
}

func (s *ApiIntegrationTestSuite) TestPostNewUser(c *C) {
	// user := Document{
	// 	"username": "pjvds",
	// 	"password": "foobar",
	// 	"email":    "pj@born2code.net:",
	// }
	//data := user.ToJson()
	//buf := bytes.NewBuffer(data)

	//_, err := http.Post("http://api.localhost:8000/users", "application/json", buf)
	b := bytes.NewBufferString("{ \"username\": \"pjvds\", \"password\": \"foobar\", \"email\": \"pj@born2code.net\" }")
	_, err := http.Post("http://api.localhost:8000/users", "application/json", b)
	if err != nil {
		c.Log("Foo")
		c.Log(err)
	}

	// c.Assert(err, IsNil)
	// c.Assert(response.StatusCode, Equals, http.StatusOK)

	// doc, err := GetBodyAsDocument(c, response)

	// c.Assert(err, IsNil)
	// c.Assert(doc["username"], Equals, user["username"])
}

func GetBodyAsDocument(c *C, response *http.Response) (Document, error) {
	var document Document
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return document, err
	}

	err = json.Unmarshal(data, &document)
	if err != nil {
		c.Logf("RAW Json: %s", string(data))
	}

	return document, err
}

func (doc Document) ToJson() []byte {
	data, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	return data
}
