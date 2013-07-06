package mvc

import (
	"net/http"
)

type RedirectResultState struct {
	permanent bool
	url       string
}

func RedirectResult(url string, permanent bool) *RedirectResultState {
	return &RedirectResultState{
		permanent: permanent,
		url:       url,
	}
}

func (r *RedirectResultState) WriteResponse(response http.ResponseWriter) {
	if r.permanent {
		response.WriteHeader(http.StatusMovedPermanently)
	} else {
		response.WriteHeader(http.StatusTemporaryRedirect)
	}
	response.Write([]byte(r.url))
}
