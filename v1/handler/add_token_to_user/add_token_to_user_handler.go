package add_token_to_user

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type addTokenToUser func(model.AuthToken, model.UserName) error

type handler struct {
	addTokenToUser addTokenToUser
}

func New(addTokenToUser addTokenToUser) *handler {
	h := new(handler)
	h.addTokenToUser = addTokenToUser
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("add token to user")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("add token to user failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("add token to user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.UsernameTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	return h.addTokenToUser(request.AuthToken, request.Userame)
}
