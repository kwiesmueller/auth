package access_denied

import (
	"net/http"

	"fmt"

	"github.com/golang/glog"
)

type handler struct {
	message string
	status  int
}

func New() *handler {
	h := new(handler)
	h.status = http.StatusForbidden
	h.message = http.StatusText(h.status)
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("%d => %s", h.status, h.message)
	http.Error(resp, fmt.Sprintf(h.message), h.status)
}
