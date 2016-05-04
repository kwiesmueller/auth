package filter

import (
	"net/http"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type handler struct {
	check   Check
	success http.HandlerFunc
	failure http.HandlerFunc
}

type Check func(req *http.Request) (bool, error)

func New(check Check, success http.HandlerFunc, failure http.HandlerFunc) *handler {
	h := new(handler)
	h.check = check
	h.success = success
	h.failure = failure
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	handlerFunc := h.handlerFunc(req)
	handlerFunc(resp, req)
}

func (h *handler) handlerFunc(req *http.Request) http.HandlerFunc {
	result, err := h.check(req)
	if err != nil {
		logger.Debugf("check failed: %v", err)
		return h.failure
	}
	if !result {
		return h.failure
	}
	return h.success
}
