package host

import (
	"fmt"
	"github.com/gorilla/handlers"
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
		githubCtlr: controllers.NewGithubOAuthController(config.Github.ClientId, config.Github.ClientSecret,
			config.Github.AuthorizeUrl, config.Github.AccessTokenUrl, userRepository),
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

	router.Handle("/", handlers.MethodHandler{
		"GET": HttpReponseWrapper(s.homeCtlr.HandleIndex),
	})

	router.Handle("/ping", handlers.MethodHandler{
		"GET": HttpReponseWrapper(s.homeCtlr.HandlePing),
	})

	router.Handle("/auth/github/authorizeurl", handlers.MethodHandler{
		"GET": HttpReponseWrapper(s.githubCtlr.GetGithubAuthorizeUrl),
	})
	router.Handle("/auth/github/callback", handlers.MethodHandler{
		"GET": HttpReponseWrapper(s.githubCtlr.GithubCallback),
	})

	addCallbackHandler := mvc.NewJsonBodyRequestArgsObjectHandler(s.callbackCtlr.NewCallback)
	router.Handle("/user/callbacks", handlers.MethodHandler{
		"GET": s.authenticator.Wrap(func(response http.ResponseWriter, request *mvc.AuthenticatedRequest) {
			s.callbackCtlr.ListCallbacks(request).WriteResponse(response)
		}),
		"POST": s.authenticator.Wrap(addCallbackHandler.ServeAuthHTTP),
	})

	loginHandler := mvc.NewJsonBodyRequestArgsObjectHandler(s.userCtlr.Login)
	router.Handle("/login", handlers.MethodHandler{
		"POST": loginHandler,
	})

	router.Handle("/user/{id}", handlers.MethodHandler{
		"GET": mvc.HandlerFuncToHandler(func(response http.ResponseWriter, req *http.Request) {
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
		}),
	})

	addUserHandler := mvc.NewJsonBodyRequestArgsObjectHandler(s.userCtlr.AddUser)
	router.Handle("/users", handlers.MethodHandler{
		"POST": addUserHandler,
	})

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
	response.Header().Add("Access-Control-Allow-Origin", "*")
	s.router.ServeHTTP(response, request)
}
