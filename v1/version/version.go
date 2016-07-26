package version

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
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
	logger.Debugf("version")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("version failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("version success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	logger.Debugf("version")
	versionNumber := v1.VERSION
	resp.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(resp).Encode(&versionNumber)
}
