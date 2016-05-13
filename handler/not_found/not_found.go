package not_found

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
	logger.Debugf("not found %s %s", req.Method, req.URL.Path)
	resp.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(resp, "not found %s %s\n", req.Method, req.URL.Path)
}
