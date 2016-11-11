package remove_token_from_user

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type removeTokenFromUser func(model.AuthToken, model.UserName) error

type handler struct {
	removeTokenFromUser removeTokenFromUser
}

func New(removeTokenFromUser removeTokenFromUser) *handler {
	h := new(handler)
	h.removeTokenFromUser = removeTokenFromUser
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("remove token to user")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("remove token to user failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("remove token to user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.UsernameTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	return h.removeTokenFromUser(request.AuthToken, request.Userame)
}
