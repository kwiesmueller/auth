package router

import (
	"net/http"

	"fmt"

	"github.com/bborbe/auth/v1"
	"github.com/bborbe/http_handler/adapter"
	"github.com/gorilla/mux"
)

type handler struct {
	router http.Handler
}

func New(
	prefix string,
	notFound http.HandlerFunc,
	check http.HandlerFunc,
	v1_router http.HandlerFunc,
) *handler {
	router := mux.NewRouter()

	router.PathPrefix(fmt.Sprintf("%s/api/%s", prefix, v1.VERSION)).HandlerFunc(v1_router)

	router.Path(fmt.Sprintf("%s/healthz", prefix)).Methods("GET").HandlerFunc(check)
	router.Path(fmt.Sprintf("%s/readiness", prefix)).Methods("GET").HandlerFunc(check)

	router.NotFoundHandler = adapter.New(notFound)

	h := new(handler)
	h.router = router
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(resp, req)
}
