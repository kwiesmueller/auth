package user_delete

import (
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type DeleteUser func(username model.UserName) error

type handler struct {
	deleteUser DeleteUser
}

func New(
	deleteUser DeleteUser,
) *handler {
	h := new(handler)
	h.deleteUser = deleteUser
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("delete user")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("delete user failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("delete user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	parts := strings.Split(req.RequestURI, "/")
	if len(parts) == 0 {
		glog.V(2).Infof("auth token missing")
		return fmt.Errorf("invalid request uri: %s", req.RequestURI)
	}
	username := model.UserName(parts[len(parts)-1])
	if err := h.deleteUser(username); err != nil {
		return err
	}
	resp.WriteHeader(http.StatusOK)
	return nil
}
