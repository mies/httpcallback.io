package api

import (
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
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
	UserId    model.ObjectId `json:"userId"`
	Username  string         `json:"username"`
	AuthToken string         `json:"authToken"`
}

func (ctr *UserController) AddUser(request *http.Request, args *AddUserRequest) (*JsonResponse, error) {
	creationDate := time.Now()
	newUser := &model.User{
		Id:           model.NewObjectId(),
		CreatedAt:    creationDate,
		Username:     args.Username,
		PasswordHash: security.HashPassword(args.Username, args.Password, creationDate),
		AuthToken:    security.NewAuthToken(),
	}

	if err := ctr.users.Add(newUser); err != nil {
		return nil, err
	}

	return JsonResult(&AddUserResponse{
		UserId:    newUser.Id,
		Username:  newUser.Username,
		AuthToken: newUser.AuthToken.String(),
	})
}

func (ctr *UserController) ListUsers(request *http.Request) (*JsonResponse, error) {
	users, err := ctr.users.List()
	if err != nil {
		return nil, err
	}

	return JsonResult(users)
}
