package user_register

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type createUserWithToken func(userName model.UserName, authToken model.AuthToken) error

type handler struct {
	createUserWithToken createUserWithToken
}

func New(addTokenToUser createUserWithToken) *handler {
	h := new(handler)
	h.createUserWithToken = addTokenToUser
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("register user")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("register user failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
		return
	}
	glog.V(3).Infof("register user success")
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.RegisterRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		glog.V(3).Infof("decode json failed: %v", err)
		return err
	}
	glog.V(4).Infof("register user %v with token %v", request.UserName, request.AuthToken)
	err := h.createUserWithToken(request.UserName, request.AuthToken)
	if err != nil {
		glog.V(3).Infof("register user %v failed: %v", request.UserName, err)
		return err
	}
	glog.V(4).Infof("register user %v successful", request.UserName)
	return nil
}
