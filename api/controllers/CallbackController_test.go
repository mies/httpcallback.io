package controllers

import (
	"github.com/pjvds/httpcallback.io/data/memory"
	. "launchpad.net/gocheck"
	"testing"
)

type CallbackControllerTestSuite struct {
	callbackRepository memory.MemoryCallbackRepository
}

func NewCallbackControllerTestSuite() *CallbackControllerTestSuite {
	return &CallbackControllerTestSuite{
		callbackRepository: memory.NewMemoryCallbackRepository(),
	}
}

func (t *CallbackControllerTestSuite) SetUpSuite(c *C) {
}
