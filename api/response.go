package api

import ()

type ResponseWriter interface {
	SetHeader(string, string)
	Success(Response)
	ErrMethodNotAllowed()
}
