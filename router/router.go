package router

import (
	"net/http"

	"github.com/bborbe/server/handler/adapter"
	"github.com/gorilla/mux"
)

type handler struct {
	router http.Handler
}

func New(
	notFound http.HandlerFunc,
	check http.HandlerFunc,
	login http.HandlerFunc,
	applicationCreate http.HandlerFunc,
	applicationDelete http.HandlerFunc,
	applicationGet http.HandlerFunc,
	userRegister http.HandlerFunc,
	userUnregister http.HandlerFunc,
	userDelete http.HandlerFunc,
	tokenAdd http.HandlerFunc,
	tokenRemove http.HandlerFunc,
	userGroupAdd http.HandlerFunc,
	userGroupRemove http.HandlerFunc,
	userDataSet http.HandlerFunc,
	userDataSetValue http.HandlerFunc,
	userDataGet http.HandlerFunc,
	userDataGetValue http.HandlerFunc,
	userDataDelete http.HandlerFunc,
	userDataDeleteValue http.HandlerFunc,
) *handler {
	router := mux.NewRouter()

	router.Path("/user/{username}/data").Methods("POST").HandlerFunc(userDataSet)
	router.Path("/user/{username}/data/{key}").Methods("POST").HandlerFunc(userDataSetValue)
	router.Path("/user/{username}/data").Methods("GET").HandlerFunc(userDataGet)
	router.Path("/user/{username}/data/{key}").Methods("GET").HandlerFunc(userDataGetValue)
	router.Path("/user/{username}/data").Methods("DELETE").HandlerFunc(userDataDelete)
	router.Path("/user/{username}/data/{key}").Methods("DELETE").HandlerFunc(userDataDeleteValue)

	router.Path("/healthz").Methods("GET").HandlerFunc(check)
	router.Path("/readiness").Methods("GET").HandlerFunc(check)

	router.Path("/login").Methods("POST").HandlerFunc(login)

	router.Path("/application").Methods("POST").HandlerFunc(applicationCreate)
	router.PathPrefix("/application/").Methods("DELETE").HandlerFunc(applicationDelete)
	router.PathPrefix("/application/").Methods("GET").HandlerFunc(applicationGet)

	router.Path("/user/{username}").Methods("DELETE").HandlerFunc(userDelete)
	router.Path("/user").Methods("POST").HandlerFunc(userRegister)

	router.Path("/token/{token}").Methods("DELETE").HandlerFunc(userUnregister)
	router.Path("/token").Methods("POST").HandlerFunc(tokenAdd)
	router.Path("/token").Methods("DELETE").HandlerFunc(tokenRemove)

	router.Path("/user_group").Methods("POST").HandlerFunc(userGroupAdd)
	router.Path("/user_group").Methods("DELETE").HandlerFunc(userGroupRemove)

	router.NotFoundHandler = adapter.New(notFound)

	h := new(handler)
	h.router = router
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(resp, req)
}
