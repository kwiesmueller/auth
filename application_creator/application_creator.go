package application_creator

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

const PASSWORD_LENGTH = 16

type CreateApplication func(application api.Application) error
type GeneratePassword func(length int) string
type handler struct {
	createApplication CreateApplication
	generatePassword  GeneratePassword
}

func New(createApplication CreateApplication, generatePassword GeneratePassword) *handler {
	h := new(handler)
	h.createApplication = createApplication
	h.generatePassword = generatePassword
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
	application := api.Application{
		ApplicationName:     request.ApplicationName,
		ApplicationPassword: api.ApplicationPassword(h.generatePassword(PASSWORD_LENGTH)),
	}
	err := h.createApplication(application)
	if err != nil {
		return err
	}
	response.ApplicationName = application.ApplicationName
	response.ApplicationPassword = application.ApplicationPassword
	return nil
}
