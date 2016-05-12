package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type handler struct {
	router http.Handler
}

func New(
	check HandlerFunc,
	login HandlerFunc,
	applicationCreate HandlerFunc,
	applicationDelete HandlerFunc,
	applicationGet HandlerFunc,
	userRegister HandlerFunc,
	userUnregister HandlerFunc,
	tokenAdd HandlerFunc,
	tokenRemove HandlerFunc,
	userGroupAdd HandlerFunc,
	userGroupRemove HandlerFunc,
	userDataSet HandlerFunc,
	userDataSetValue HandlerFunc,
	userDataGet HandlerFunc,
	userDataGetValue HandlerFunc,
	userDataDelete HandlerFunc,
	userDataDeleteValue HandlerFunc,
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

	router.Path("/user").Methods("POST").HandlerFunc(userRegister)
	router.PathPrefix("/user/").Methods("DELETE").HandlerFunc(userUnregister)

	router.Path("/token").Methods("POST").HandlerFunc(tokenAdd)
	router.Path("/token").Methods("DELETE").HandlerFunc(tokenRemove)

	router.Path("/user_group").Methods("POST").HandlerFunc(userGroupAdd)
	router.Path("/user_group").Methods("DELETE").HandlerFunc(userGroupRemove)

	h := new(handler)
	h.router = router
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(resp, req)
}
