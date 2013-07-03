package api

import (
	"github.com/pjvds/httpcallback.io/model"
	. "launchpad.net/gocheck"
	"net/http"
)

// Setup the test suite
var _ = Suite(&AuthenticatedRequestTestSuite{})

type AuthenticatedRequestTestSuite struct{}

func (s *AuthenticatedRequestTestSuite) SetUpSuite(c *C) {
}

func (s *AuthenticatedRequestTestSuite) TestAuthenticatedRequestCtorSetsCorrectFields(c *C) {
	request := &http.Request{}
	userId := model.NewObjectId()
	username := "username"

	sut := NewAuthenticatedRequest(request, userId, username)
	c.Assert(sut.Request, Equals, request)
	c.Assert(sut.UserId, Equals, userId)
	c.Assert(sut.Username, Equals, username)
}

func (s *AuthenticatedRequestTestSuite) TestAuthenticatedRequestCtorPanicsOnNilRequest(c *C) {
	request := (*http.Request)(nil)
	userId := model.NewObjectId()
	username := "username"

	c.Assert(func() {
		NewAuthenticatedRequest(request, userId, username)
	}, Panics, "request cannot be nil")
}

func (s *AuthenticatedRequestTestSuite) TestAuthenticatedRequestCtorPanicsOnEmptyUsername(c *C) {
	request := &http.Request{}
	userId := model.NewObjectId()
	var username string

	c.Assert(func() {
		NewAuthenticatedRequest(request, userId, username)
	}, Panics, "username cannot be empty")
}
