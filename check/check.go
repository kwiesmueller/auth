package check

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
	if check() {
		logger.Debugf("check => ok")
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, "ok")
	} else {
		logger.Debugf("check => failed")
		http.Error(resp, fmt.Sprintf("check failed"), http.StatusInternalServerError)
	}
}

func check() bool {
	return true
}
