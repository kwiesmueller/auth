package token_by_username

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type listTokenOfUser func(model.UserName) ([]model.AuthToken, error)

type handler struct {
	listTokenOfUser listTokenOfUser
}

func New(removeTokenToUserWithToken listTokenOfUser) *handler {
	h := new(handler)
	h.listTokenOfUser = removeTokenToUserWithToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("list user")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("list user failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("list user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	username := model.UserName(req.FormValue("username"))
	var err error
	var tokens []model.AuthToken
	if tokens, err = h.listTokenOfUser(username); err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&tokens)
}
