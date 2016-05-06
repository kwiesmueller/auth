package application_getter

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type GetApplication func(applicationName api.ApplicationName) (*api.Application, error)
type IsApplicationNotFound func(err error) bool

type handler struct {
	getApplication        GetApplication
	isApplicationNotFound IsApplicationNotFound
}

func New(getApplication GetApplication, isApplicationNotFound IsApplicationNotFound) *handler {
	h := new(handler)
	h.getApplication = getApplication
	h.isApplicationNotFound = isApplicationNotFound
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("get application")
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
	application, err := h.getApplication(api.ApplicationName(last))
	if err != nil {
		if h.isApplicationNotFound(err) {
			e := error_handler.NewErrorMessage(http.StatusNotFound, err.Error())
			e.ServeHTTP(resp, req)
			return nil
		}
		return err
	}
	return json.NewEncoder(resp).Encode(&api.GetApplicationResponse{
		ApplicationName:     application.ApplicationName,
		ApplicationPassword: application.ApplicationPassword,
	})
}
