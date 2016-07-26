package application_getter

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type GetApplication func(applicationName model.ApplicationName) (*model.Application, error)

type handler struct {
	getApplication GetApplication
}

func New(getApplication GetApplication) *handler {
	h := new(handler)
	h.getApplication = getApplication
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
	application, err := h.getApplication(model.ApplicationName(last))
	if err != nil {
		e := error_handler.NewErrorMessage(http.StatusNotFound, err.Error())
		e.ServeHTTP(resp, req)
		return nil
	}
	return json.NewEncoder(resp).Encode(&v1.GetApplicationResponse{
		ApplicationName:     application.ApplicationName,
		ApplicationPassword: application.ApplicationPassword,
	})
}