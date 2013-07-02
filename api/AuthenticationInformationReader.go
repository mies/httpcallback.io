package api

import (
	"net/http"
	"regexp"
)

// Reads authentication information from http request. It reads it from the
// header, or if not present, from the query string. If that also isn't present
// empty values are returned.
func GetAuthorizationInfoFromRequest(request *http.Request) (username string, token string) {
	authHeader := request.Header.Get("Authorization")
	if authHeader != "" {
		Log.Debug("Received request with authorization header: %v", authHeader)

		pattern := "(?P<type>HttpCallbackLogin).*(username=\")(?P<username>\\w+).*(token=\")(?P<token>\\w+)"
		regex := regexp.MustCompile(pattern)
		match := regex.FindStringSubmatch(authHeader)

		if match == nil {
			return
		}

		var authType string
		for i, name := range regex.SubexpNames() {
			// Ignore the whole regexp match and unnamed groups
			if i == 0 || name == "" {
				continue
			}

			switch name {
			case "type":
				authType = match[i]
			case "username":
				username = match[i]
			case "token":
				token = match[i]
			default:
				Log.Warning("unknown match in authentication header:\n\tname: %v\n\tvalue: %v", name, match[i])
			}

			if authType == "" {
				Log.Warning("No auth type specified in header value, clearing username and token")
				username = ""
				token = ""
			}
		}
	} else {
		username = request.URL.Query().Get("auth_username")
		token = request.URL.Query().Get("auth_token")
	}

	return
}
