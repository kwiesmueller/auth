package delete

import (
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type deleteApplication func(applicationName model.ApplicationName) error

type handler struct {
	deleteApplication deleteApplication
}

func New(deleteApplication deleteApplication) *handler {
	h := new(handler)
	h.deleteApplication = deleteApplication
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("delete application started")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("delete application failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
		return
	}
	glog.V(3).Infof("delete application finished")
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
	glog.V(4).Infof("application deleted")
	return nil
}
