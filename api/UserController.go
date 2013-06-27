package api

import (
	"fmt"
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
	Id        model.ObjectId            `json:"id"`
	Username  string                    `json:"username"`
	AuthToken model.AuthenticationToken `json:"authToken"`
}

func (ctr *UserController) AddUser(request *http.Request, args *AddUserRequest) (HttpResponse, error) {
	Log.Info("Handling AddUser request for new user with username: %s", args.Username)

	creationDate := time.Now()

	newUser := model.User{
		Id:           model.NewObjectId(),
		CreatedAt:    creationDate,
		Username:     args.Username,
		PasswordHash: security.HashPassword(args.Username, args.Password, creationDate),
		AuthToken:    security.NewAuthToken(),
	}

	if err := ctr.users.Add(&newUser); err != nil {
		fmt.Println("Unable to add new user. Error from repository:,", err)
		return nil, err
	}

	return JsonResult(AddUserResponse{
		Id:        newUser.Id,
		Username:  newUser.Username,
		AuthToken: newUser.AuthToken,
	})
}

type GetUserRequestArgs struct {
	UserId string
}

func (ctr *UserController) GetUser(request *http.Request, args *GetUserRequestArgs) (HttpResponse, error) {
	userId, err := model.ParseObjectId(args.UserId)
	if err != nil {
		fmt.Errorf("Invalid user id '%s': %s\nWill return 404 to user.", args.UserId, err.Error())
		return NewHttpStatusCodeResult(http.StatusNotFound), nil
	}

	user, err := ctr.users.Get(userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		Log.Debug("No user found with id '%s', returning 404 not found.", args.UserId)
		return NewHttpStatusCodeResult(http.StatusNotFound), nil
	}

	return JsonResult(struct {
		Id        model.ObjectId `json:"id"`
		Username  string         `json:"username"`
		CreatedAt time.Time      `json:"createdAt"`
	}{Id: user.Id,
		Username:  user.Username,
		CreatedAt: user.CreatedAt})
}
