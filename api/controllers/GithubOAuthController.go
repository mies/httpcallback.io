package controllers

import (
	"bitbucket.org/gosimple/oauth2"
	"encoding/json"
	"errors"
	"github.com/pjvds/httpcallback.io/data"
	. "github.com/pjvds/httpcallback.io/mvc"
	"net/http"
)

const (
	// The base url of the github API which also include the oauth2
	// methods. Which are also the only methods we will use in this controller
	BaseUrl = "https://api.github.com/"
)

type GithubOAuthController struct {
	GithubService  *oauth2.OAuth2Service
	UserRepository data.UserRepository
}

type GithubAccessTokenRequest struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func NewGithubOAuthController(clientId, clientSecret, authorizeUrl, accessTokenUrl string, userRepository data.UserRepository) *GithubOAuthController {
	return &GithubOAuthController{
		GithubService:  oauth2.Service(clientId, clientSecret, authorizeUrl, accessTokenUrl),
		UserRepository: userRepository,
	}
}

func (ctr *GithubOAuthController) GithubCallback(request *http.Request) ActionResult {
	code := request.URL.Query().Get("code")
	Log.Debug("Received callback from github oauth provider with code: %v", code)

	service := ctr.GithubService

	// Get access token.
	token, err := service.GetAccessToken(code)
	if err != nil {
		Log.Error("Get access token error: ", err)
		return ErrorResult(err)
	}
	Log.Debug("Received access token: %v", token)

	// Prepare resource request.
	github := oauth2.Request(BaseUrl, token.AccessToken)
	github.Header.Add("Accept", "application/vnd.github.v3")
	github.AccessTokenInHeader = true
	github.AccessTokenInHeaderScheme = "token"

	// Make the request.
	// Provide API end point (http://developer.github.com/v3/users/#get-the-authenticated-user)
	apiEndPoint := "user/emails"
	userEmailsResponse, err := github.Get(apiEndPoint)
	if err != nil {
		Log.Error("Get: ", err)
		return ErrorResult(err)
	}
	defer userEmailsResponse.Body.Close()

	var docs []JsonDocument
	var decoder = json.NewDecoder(userEmailsResponse.Body)
	if err := decoder.Decode(&docs); err != nil {
		Log.Error("Error while decoding github response body: %v", err)
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
	authUrl := ctr.GithubService.GetAuthorizeURL("")

	return JsonResult(JsonDocument{
		"authorizeUrl": authUrl,
	})
}
