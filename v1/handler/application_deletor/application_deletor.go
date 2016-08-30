package application_deletor

import (
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type DeleteApplication func(applicationName model.ApplicationName) error

type handler struct {
	deleteApplication DeleteApplication
}

func New(deleteApplication DeleteApplication) *handler {
	h := new(handler)
	h.deleteApplication = deleteApplication
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
	parts := strings.Split(req.RequestURI, "/")
	if len(parts) == 0 {
		return fmt.Errorf("invalid request uri: %s", req.RequestURI)
	}
	last := parts[len(parts)-1]

	err := h.deleteApplication(model.ApplicationName(last))
	if err != nil {
		return err
	}
	glog.V(2).Infof("application deleted")
	return nil
}
