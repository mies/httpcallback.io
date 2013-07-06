package controllers

import (
	"bitbucket.org/gosimple/oauth2"
	"encoding/json"
	"errors"
	. "github.com/pjvds/httpcallback.io/mvc"
	"io/ioutil"
	"net/http"
)

type GithubOAuthController struct {
	ClientId       string
	ClientSecret   string
	AuthorizeUrl   string
	AccessTokenUrl string
}

type GithubAccessTokenRequest struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func NewGithubOAuthController(clientId, clientSecret, authorizeUrl, accessTokenUrl string) *GithubOAuthController {
	return &GithubOAuthController{
		ClientId:       clientId,
		ClientSecret:   clientSecret,
		AuthorizeUrl:   authorizeUrl,
		AccessTokenUrl: accessTokenUrl,
	}
}

func (ctr *GithubOAuthController) GithubCallback(request *http.Request) ActionResult {
	code := request.URL.Query().Get("code")
	Log.Debug("Received callback from github oauth provider with code: %v", code)

	service := oauth2.Service(
		ctr.ClientId,
		ctr.ClientSecret,
		ctr.AuthorizeUrl,
		ctr.AccessTokenUrl)

	// Get access token.
	token, err := service.GetAccessToken(code)
	if err != nil {
		Log.Error("Get access token error: ", err)
		return ErrorResult(err)
	}
	Log.Debug("Received access token: %v", token)

	// Prepare resource request.
	apiBaseURL := "https://api.github.com/"
	github := oauth2.Request(apiBaseURL, token.AccessToken)
	github.Header.Add("Accept", "application/vnd.github.v3")
	github.AccessTokenInHeader = true
	github.AccessTokenInHeaderScheme = "token"

	// Make the request.
	// Provide API end point (http://developer.github.com/v3/users/#get-the-authenticated-user)
	apiEndPoint := "user/emails"
	githubUserData, err := github.Get(apiEndPoint)
	if err != nil {
		Log.Error("Get: ", err)
		return ErrorResult(err)
	}
	defer githubUserData.Body.Close()

	body, err := ioutil.ReadAll(githubUserData.Body)
	if err != nil {
		Log.Error("Error while reading body from response: ", err)
		return ErrorResult(err)
	}

	var docs []map[string]interface{}
	if err := json.Unmarshal(body, &docs); err != nil {
		Log.Error("Error while decoding github response body: %v\n\n\tBody: %v", err, string(body))
		return ErrorResult(err)
	}

	for _, doc := range docs {
		if doc["primary"].(bool) == true {
			return JsonResult(&JsonDocument{
				"email":    doc["email"],
				"verified": doc["verified"],
			})
		}
	}

	return ErrorResult(errors.New("No primary email address found in github response"))
}

func (ctr *GithubOAuthController) GetGithubAuthorizeUrl(request *http.Request) ActionResult {
	service := oauth2.Service(
		ctr.ClientId,
		ctr.ClientSecret,
		ctr.AuthorizeUrl,
		ctr.AccessTokenUrl)
	service.Scope = "user:email"

	// Get authorization url.
	authUrl := service.GetAuthorizeURL("")

	return JsonResult(JsonDocument{
		"authorizeUrl": authUrl,
	})
}
