package controllers

import (
	"bitbucket.org/gosimple/oauth2"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/pjvds/httpcallback.io/api/messages"
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
	. "github.com/pjvds/httpcallback.io/mvc"
	"github.com/pjvds/httpcallback.io/security"
	"io/ioutil"
	"net/http"
	"time"
)

type UserController struct {
	users data.UserRepository
}

func NewUserController(users data.UserRepository) *UserController {
	return &UserController{
		users: users,
	}
}

type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AddUserResponse struct {
	Id        model.ObjectId            `bson:"_id,omitempty" json:"id"`
	Username  string                    `json:"username"`
	AuthToken model.AuthenticationToken `json:"authToken"`
}

func (ctr *UserController) GetGithubAuthorizeUrl(request *http.Request) ActionResult {
	service := oauth2.Service(
		"ce1b9a918b75be5da302",
		"50f90c49c516c065ad322518e9b43d5cd9936070",
		"https://github.com/login/oauth/authorize",
		"https://github.com/login/oauth/access_token")
	service.Scope = "user:email"

	// Get authorization url.
	authUrl := service.GetAuthorizeURL("")

	return JsonResult(JsonDocument{
		"authorizeUrl": authUrl,
	})
}

type GithubAccessTokenRequest struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (ctr *UserController) GithubCallback(request *http.Request) ActionResult {
	code := request.URL.Query().Get("code")
	Log.Debug("Received callback from github oauth provider with code: %v", code)

	service := oauth2.Service(
		"ce1b9a918b75be5da302",
		"50f90c49c516c065ad322518e9b43d5cd9936070",
		"https://github.com/login/oauth/authorize",
		"https://github.com/login/oauth/access_token")

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

func (ctr *UserController) AddUser(request *http.Request, args *AddUserRequest) ActionResult {
	log.Info("Handling AddUser request for new user with username: %v", args.Username)

	creationDate := time.Now()

	newUser := model.User{
		Id:           model.NewObjectId(),
		CreatedAt:    creationDate,
		Username:     args.Username,
		PasswordHash: security.HashPassword(args.Username, args.Password, creationDate),
		AuthToken:    security.NewAuthToken(),
	}

	if err := ctr.users.Add(&newUser); err != nil {
		Log.Error("Unable to add new user. Error from repository:,", err)
		return ErrorResult(err)
	}

	return JsonResult(AddUserResponse{
		Id:        newUser.Id,
		Username:  newUser.Username,
		AuthToken: newUser.AuthToken,
	})
}

func (ctr *UserController) GetUser(request *http.Request, args *GetUserRequest) ActionResult {
	userId, err := model.ParseObjectId(args.UserId)
	if err != nil {
		// TODO: Invalid request!
		return NotFoundResult(fmt.Sprintf("user with id '%v' does not exist", args.UserId))
	}

	user, err := ctr.users.Get(userId)
	if err != nil {
		return ErrorResult(err)
	}

	if user == nil {
		return NotFoundResult(fmt.Sprintf("user with id '%v' does not exist", args.UserId))
	}

	return JsonResult(struct {
		Id        model.ObjectId `bson:"_id,omitempty" json:"id"`
		Username  string         `json:"username"`
		CreatedAt time.Time      `json:"createdAt"`
	}{Id: user.Id,
		Username:  user.Username,
		CreatedAt: user.CreatedAt})
}
