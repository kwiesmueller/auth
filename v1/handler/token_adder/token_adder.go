package token_adder

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type addTokenToUserWithToken func(newToken model.AuthToken, userToken model.AuthToken) error

type handler struct {
	addTokenToUserWithToken addTokenToUserWithToken
}

func New(addTokenToUserWithToken addTokenToUserWithToken) *handler {
	h := new(handler)
	h.addTokenToUserWithToken = addTokenToUserWithToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("add token")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("add token failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("add token success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.AddTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	glog.V(2).Infof("add token %v to user with token %v", request.Token, request.AuthToken)
	if err := h.addTokenToUserWithToken(request.Token, request.AuthToken); err != nil {
		return err
	}
	return nil
}
