package token_remover

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type RemoveTokenToUserWithToken func(newToken model.AuthToken, userToken model.AuthToken) error

type handler struct {
	removeTokenToUserWithToken RemoveTokenToUserWithToken
}

func New(removeTokenToUserWithToken RemoveTokenToUserWithToken) *handler {
	h := new(handler)
	h.removeTokenToUserWithToken = removeTokenToUserWithToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("remove token")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("remove token failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
		return
	}
	glog.V(2).Infof("remove token success")
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.RemoveTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.removeTokenToUserWithToken(request.Token, request.AuthToken)
	if err != nil {
		return err
	}
	return nil
}
