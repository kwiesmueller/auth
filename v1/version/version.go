package version

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type handler struct {
}

func New() *handler {
	h := new(handler)
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("version")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("version failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("version success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	glog.V(2).Infof("version")
	versionNumber := v1.VERSION
	resp.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(resp).Encode(&versionNumber)
}
