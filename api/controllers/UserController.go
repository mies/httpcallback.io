package controllers

import (
	"fmt"
	. "github.com/pjvds/httpcallback.io/api/messages"
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
	. "github.com/pjvds/httpcallback.io/mvc"
	"github.com/pjvds/httpcallback.io/security"
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// func (ctr *UserController) Login(request *http.Request, args *LoginRequest) ActionResult {
// 	hash := security.HashPassword(args.Username, args.Password)
// 	user, err := ctr.users.GetByUsernameAndPasswordHash(args.Username, hash)
// 	if err != nil {
// 		return ErrorResult(err)
// 	}
// 	if user == nil {
// 		return NotFoundResult("no user found by provided credentials")
// 	}
// 	return JsonResult(JsonDocument{
// 		"token": user.AuthToken.String(),
// 	})
// }

func (ctr *UserController) AddUser(request *http.Request, args *AddUserRequest) ActionResult {
	log.Info("Handling AddUser request for new user with username: %v", args.Username)

	creationDate := time.Now()

	newUser := model.NewUser(
		model.NewObjectId(),
		creationDate,
		args.Username,
		security.HashPassword(args.Username, args.Password),
		args.Email)
	newUser.AddAuthToken(security.NewAuthToken())

	if err := ctr.users.Add(newUser); err != nil {
		Log.Error("Unable to add new user. Error from repository:,", err)
		return ErrorResult(err)
	}

	return JsonResult(AddUserResponse{
		Id:        model.ObjectId(newUser.Id()),
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
	}{Id: model.ObjectId(user.Id()),
		Username:  user.Username,
		CreatedAt: user.CreatedAt})
}
