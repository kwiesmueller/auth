package create

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type createApplication func(applicationName model.ApplicationName) (*model.Application, error)

type handler struct {
	createApplication createApplication
}

func New(createApplication createApplication) *handler {
	h := new(handler)
	h.createApplication = createApplication
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("create application started")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("create application failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
		return
	}
	glog.V(3).Infof("create application finished")
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request model.Application
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		glog.V(3).Infof("parse request failed: %v", err)
		return err
	}
	application, err := h.createApplication(request.ApplicationName)
	if err != nil {
		glog.V(3).Infof("create application failed: %v", err)
		return err
	}
	return json.NewEncoder(resp).Encode(&application)
}
