package access_denied

import (
	"net/http"

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
	logger.Debugf("unauthorized")
	resp.WriteHeader(http.StatusUnauthorized)
}
