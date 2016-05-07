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
	userCreate HandlerFunc,
) *handler {
	router := mux.NewRouter()
	router.Path("/healthz").Methods("GET").HandlerFunc(check)
	router.Path("/readiness").Methods("GET").HandlerFunc(check)
	router.Path("/login").Methods("POST").HandlerFunc(login)
	router.Path("/application").Methods("POST").HandlerFunc(applicationCreate)
	router.PathPrefix("/application/").Methods("DELETE").HandlerFunc(applicationDelete)
	router.PathPrefix("/application/").Methods("GET").HandlerFunc(applicationGet)
	router.Path("/user").Methods("POST").HandlerFunc(userCreate)

	h := new(handler)
	h.router = router
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(resp, req)
}
