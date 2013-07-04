package host

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pjvds/httpcallback.io/api"
	"github.com/pjvds/httpcallback.io/api/controllers"
	"github.com/pjvds/httpcallback.io/api/messages"
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/mvc"
	"net/http"
)

type HttpCallbackApiServer struct {
	router *mux.Router
}

func NewServer(repositoryFactory data.RepositoryFactory) *HttpCallbackApiServer {
	callbacksController := controllers.NewCallbackController(repositoryFactory.CreateCallbackRepository())
	usersController := controllers.NewUserController(repositoryFactory.CreateUserRepository())
	homeController := controllers.NewHomeController()

	authenticator := mvc.NewAuthenticator(repositoryFactory.CreateUserRepository())
	router := mux.NewRouter()

	apiPostRouter := router.Methods("POST").Subrouter()
	apiGetRouter := router.Methods("GET").Subrouter()
	apiGetRouter.HandleFunc("/", HttpReponseWrapper(homeController.HandleIndex))
	apiGetRouter.HandleFunc("/ping", HttpReponseWrapper(homeController.HandlePing))
	apiGetRouter.Handle("/user/callbacks", authenticator.Wrap(func(response http.ResponseWriter, req *mvc.AuthenticatedRequest) {
		result := callbacksController.ListCallbacks(req)
		result.WriteResponse(response)
	}))
	apiGetRouter.HandleFunc("/user/{id}", func(response http.ResponseWriter, req *http.Request) {
		var result mvc.ActionResult

		userId, ok := mux.Vars(req)["id"]
		if !ok {
			Log.Warning("id parameter not given, return 404 not found.")
			result = api.NewHttpStatusCodeResult(http.StatusNotFound)
		} else {
			requestArgs := &messages.GetUserRequestArgs{
				UserId: userId,
			}

			Log.Debug("Handing request to GetUser with %+v", requestArgs)
			result = usersController.GetUser(req, requestArgs)
		}

		result.WriteResponse(response)
	})

	addUserHandler := mvc.NewJsonBodyRequestArgsObjectHandler(usersController.AddUser)
	apiPostRouter.HandleFunc("/users", func(response http.ResponseWriter, req *http.Request) {
		addUserHandler.ServeHTTP(response, req)
	})

	addCallbackHandler := mvc.NewJsonBodyRequestArgsObjectHandler(callbacksController.NewCallback)

	apiPostRouter.Handle("/user/callbacks", authenticator.Wrap(func(response http.ResponseWriter, req *mvc.AuthenticatedRequest) {
		addCallbackHandler.ServeAuthHTTP(response, req)
	}))

	return &HttpCallbackApiServer{
		router: router,
	}
}

func HttpReponseWrapper(handler func(*http.Request) mvc.ActionResult) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("[%v] %v\n", req.Method, req.URL)
		result := handler(req)
		result.WriteResponse(res)
	}
}

func (s *HttpCallbackApiServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(response, request)
}
