package get

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type getApplication func(applicationName model.ApplicationName) (*model.Application, error)

type handler struct {
	getApplication getApplication
}

func New(getApplication getApplication) *handler {
	h := new(handler)
	h.getApplication = getApplication
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("get application started")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("Marshal json failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
		return
	}
	glog.V(3).Infof("get application finished")
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	parts := strings.Split(req.RequestURI, "/")
	if len(parts) == 0 {
		return fmt.Errorf("invalid request uri: %s", req.RequestURI)
	}
	last := parts[len(parts)-1]
	application, err := h.getApplication(model.ApplicationName(last))
	if err != nil {
		e := error_handler.NewMessage(http.StatusNotFound, err.Error())
		e.ServeHTTP(resp, req)
		return nil
	}
	return json.NewEncoder(resp).Encode(&application)
}
