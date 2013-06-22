package api

import (
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/model"
	"net/http"
)

type UserController struct {
	users data.UserRepository
}

func NewUserController(users data.UserRepository) *UserController {
	return &UserController{
		users: users,
	}
}

func (ctr *UserController) AddUser(user *model.User) {
	ctr.users.Add(user)
}

func (ctr *UserController) ListUsers(request *http.Request) (*JsonResponse, error) {
	users, err := ctr.users.List()
	if err != nil {
		return nil, err
	}

	return JsonResult(users)
}
