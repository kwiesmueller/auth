package application_creator

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type CreateApplication func(applicationName api.ApplicationName) (*api.Application, error)
type handler struct {
	createApplication CreateApplication
}

func New(createApplication CreateApplication) *handler {
	h := new(handler)
	h.createApplication = createApplication
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("create application")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("create application failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("create application success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.CreateApplicationRequest
	var response api.CreateApplicationResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *api.CreateApplicationRequest, response *api.CreateApplicationResponse) error {
	application, err := h.createApplication(request.ApplicationName)
	if err != nil {
		return err
	}
	response.ApplicationName = application.ApplicationName
	response.ApplicationPassword = application.ApplicationPassword
	return nil
}
