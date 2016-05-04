package application_creator

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type CreateApplication func(application api.Application) error

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
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var createApplicationRequest api.CreateApplicationRequest
	if err := json.NewDecoder(req.Body).Decode(&createApplicationRequest); err != nil {
		return err
	}

	err := h.createApplication(api.Application{
		ApplicationName: createApplicationRequest.ApplicationName,
		ApplicationPassword: createPassword(),
	})
	if err != nil {
		return err
	}

	return json.NewEncoder(resp).Encode(&api.CreateApplicationResponse{
		ApplicationName: createApplicationRequest.ApplicationName,
		ApplicationPassword: createPassword(),
	})
}

func createPassword() api.ApplicationPassword {
	return api.ApplicationPassword("")
}
