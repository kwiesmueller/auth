package add_application

import (
	"net/http"

	"fmt"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type handler struct {
}

func New() *handler {
	h := new(handler)
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("add application")
	resp.WriteHeader(http.StatusOK)
	fmt.Fprintf(resp, "ok")
}
