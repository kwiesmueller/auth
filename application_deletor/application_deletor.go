package application_deletor

import (
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type DeleteApplication func(applicationName api.ApplicationName) error

type handler struct {
	deleteApplication DeleteApplication
}

func New(deleteApplication DeleteApplication) *handler {
	h := new(handler)
	h.deleteApplication = deleteApplication
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("delete application")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	parts := strings.Split(req.RequestURI, "/")
	if len(parts) == 0 {
		return fmt.Errorf("invalid request uri: %s", req.RequestURI)
	}
	last := parts[len(parts)-1]
	err := h.deleteApplication(api.ApplicationName(last))
	if err != nil {
		return err
	}
	logger.Debugf("application deleted")
	resp.WriteHeader(http.StatusOK)
	return nil
}
