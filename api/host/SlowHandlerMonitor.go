package host

import (
	"net/http"
	"time"
)

type SlowHandlerMonitor struct {
	threshold time.Duration
	handler   http.Handler
}

func (s *SlowHandlerMonitor) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	s.handler.ServeHTTP(response, request)
	endTime := time.Now()

	duration := endTime.Sub(startTime)
	Log.Warning("Slow request handling detected for url %v: %v", request.URL, duration.String())
}
