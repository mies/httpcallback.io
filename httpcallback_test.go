package main

import (
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

var _ = Suite(&ApiIntegrationTestSuite{
	ProcessFilename: "httpcallback.io",
})

type ApiIntegrationTestSuite struct {
	ProcessFilename string
	process         *os.Process
}

func (s *ApiIntegrationTestSuite) SetUpSuite(c *C) {
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	process, err := os.StartProcess(s.ProcessFilename, []string{"-config config.toml -port 8000"}, &procAttr)
	if err != nil {
		c.Errorf("Unable to start %s: %s", s.ProcessFilename, err.Error())
		c.Fail()
	}

	s.process = process
	time.Sleep(2 * time.Second)

	c.Logf("Started %s, pid %v", s.ProcessFilename, process.Pid)
}

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

func GetBodyAsDocument(c *C, response *http.Response) (map[string]interface{}, error) {
	var document map[string]interface{}
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
