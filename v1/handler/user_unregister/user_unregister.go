package user_unregister

import (
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type DeleteUserWithToken func(authToken model.AuthToken) error

type handler struct {
	deleteUserWithToken DeleteUserWithToken
}

func New(
	deleteUserWithToken DeleteUserWithToken,
) *handler {
	h := new(handler)
	h.deleteUserWithToken = deleteUserWithToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("unregister user")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("unregister user failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("unregister user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	parts := strings.Split(req.RequestURI, "/")
	if len(parts) == 0 {
		glog.V(2).Infof("auth token missing")
		return fmt.Errorf("invalid request uri: %s", req.RequestURI)
	}
	authToken := model.AuthToken(parts[len(parts)-1])
	if err := h.deleteUserWithToken(authToken); err != nil {
		return err
	}
	resp.WriteHeader(http.StatusOK)
	return nil
}
