package host

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pjvds/httpcallback.io/api/controllers"
	"github.com/pjvds/httpcallback.io/api/messages"
	"github.com/pjvds/httpcallback.io/data"
	"github.com/pjvds/httpcallback.io/data/memory"
	"github.com/pjvds/httpcallback.io/data/mongo"
	"github.com/pjvds/httpcallback.io/mvc"
	"net/http"
)

type HttpCallbackApiServer struct {
	router *mux.Router

	authenticator *mvc.Authenticator

	homeCtlr     *controllers.HomeController
	callbackCtlr *controllers.CallbackController
	userCtlr     *controllers.UserController
	githubCtlr   *controllers.GithubOAuthController
}

func NewServer(config *Configuration) *HttpCallbackApiServer {
	repositoryFactory, err := createRepositoryFactory(config)
	if err != nil {
		Log.Fatal("[FATAL] Could not create repository factory: " + err.Error())
	}

	userRepository := repositoryFactory.CreateUserRepository()

	server := &HttpCallbackApiServer{
		authenticator: mvc.NewAuthenticator(userRepository),
		callbackCtlr:  controllers.NewCallbackController(repositoryFactory.CreateCallbackRepository()),
		userCtlr:      controllers.NewUserController(userRepository),
		homeCtlr:      controllers.NewHomeController(),
		githubCtlr:    controllers.NewGithubOAuthController(config.Github.ClientId, config.Github.ClientSecret, config.Github.AuthorizeUrl, config.Github.AccessTokenUrl),
	}

	server.router = server.createRouter()
	return server
}

func createRepositoryFactory(config *Configuration) (data.RepositoryFactory, error) {
	if config.Mongo.UseMongo {
		Log.Debug("Running with mongo data store")
		Log.Debug("Connecting to mongo database %v", config.Mongo.DatabaseName)
		mongoSession, err := mongo.Open(config.Mongo.ServerUrl, config.Mongo.DatabaseName)
		if err != nil {
			Log.Error("Unable to connect to mongo:", err)
			return nil, err
		}
		Log.Debug("Connected succesfully")
		return mongo.NewMgoRepositoryFactory(mongoSession), nil

	} else {
		Log.Debug("Runnig with inmemory data store")
		return memory.NewMemRepositoryFactory(), nil
	}
}

func (s *HttpCallbackApiServer) createRouter() *mux.Router {
	router := mux.NewRouter()
	postRouter := router.Methods("POST").Subrouter()
	getRouter := router.Methods("GET").Subrouter()

	getRouter.HandleFunc("/", HttpReponseWrapper(s.homeCtlr.HandleIndex))
	getRouter.HandleFunc("/ping", HttpReponseWrapper(s.homeCtlr.HandlePing))

	getRouter.HandleFunc("/auth/github/authorizeurl", HttpReponseWrapper(s.githubCtlr.GetGithubAuthorizeUrl))
	getRouter.HandleFunc("/auth/github/callback", HttpReponseWrapper(s.githubCtlr.GithubCallback))

	getRouter.Handle("/user/callbacks", s.authenticator.Wrap(func(response http.ResponseWriter, request *mvc.AuthenticatedRequest) {
		s.callbackCtlr.ListCallbacks(request).WriteResponse(response)
	}))
	getRouter.HandleFunc("/user/{id}", func(response http.ResponseWriter, req *http.Request) {
		var result mvc.ActionResult

		userId, ok := mux.Vars(req)["id"]
		if !ok {
			// TODO: Invalid request!
			Log.Warning("id parameter not given, return 404 not found.")
			result = mvc.NotFoundResult("no user found with empty id")
		} else {
			requestArgs := &messages.GetUserRequest{
				UserId: userId,
			}

			Log.Debug("Handing request to GetUser with %+v", requestArgs)
			result = s.userCtlr.GetUser(req, requestArgs)
		}

		result.WriteResponse(response)
	})

	addUserHandler := mvc.NewJsonBodyRequestArgsObjectHandler(s.userCtlr.AddUser)
	postRouter.HandleFunc("/users", func(response http.ResponseWriter, req *http.Request) {
		addUserHandler.ServeHTTP(response, req)
	})

	addCallbackHandler := mvc.NewJsonBodyRequestArgsObjectHandler(s.callbackCtlr.NewCallback)
	postRouter.Handle("/user/callbacks", s.authenticator.Wrap(addCallbackHandler.ServeAuthHTTP))

	return router
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
