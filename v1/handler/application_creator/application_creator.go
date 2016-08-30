package application_creator

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type CreateApplication func(applicationName model.ApplicationName) (*model.Application, error)
type handler struct {
	createApplication CreateApplication
}

func New(createApplication CreateApplication) *handler {
	h := new(handler)
	h.createApplication = createApplication
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("create application")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("create application failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("create application success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.CreateApplicationRequest
	var response v1.CreateApplicationResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *v1.CreateApplicationRequest, response *v1.CreateApplicationResponse) error {
	application, err := h.createApplication(request.ApplicationName)
	if err != nil {
		return err
	}
	response.ApplicationName = application.ApplicationName
	response.ApplicationPassword = application.ApplicationPassword
	return nil
}
