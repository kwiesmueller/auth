package check

import (
	"net/http"

	"fmt"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Check func() error

type handler struct {
	check Check
}

func New(c Check) *handler {
	h := new(handler)
	h.check = c
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if err := h.check(); err != nil {
		logger.Debugf("check => failed: %v", err)
		http.Error(resp, fmt.Sprintf("check failed"), http.StatusInternalServerError)
	} else {
		logger.Debugf("check => ok")
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, "ok")
	}
}
