package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type handler struct {
	router http.Handler
}

func New(check HandlerFunc, login HandlerFunc, applicationCreate HandlerFunc, applicationDelete HandlerFunc) *handler {
	router := mux.NewRouter()
	router.Path("/healthz").Methods("GET").HandlerFunc(check)
	router.Path("/readiness").Methods("GET").HandlerFunc(check)
	router.Path("/login").Methods("POST").HandlerFunc(login)
	router.Path("/application").Methods("POST").HandlerFunc(applicationCreate)
	router.Path("/application").Methods("DELETE").HandlerFunc(applicationDelete)

	h := new(handler)
	h.router = router
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(resp, req)
}
