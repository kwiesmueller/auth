package verify_group_handler

import (
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/http/header"
	"github.com/golang/glog"
)

type Auth func(authToken model.AuthToken, requiredGroups []model.GroupName) (*model.UserName, error)

type handler struct {
	auth           Auth
	handler        http.Handler
	requiredGroups []model.GroupName
}

func New(subhandler http.Handler, auth Auth, requiredGroups ...model.GroupName) *handler {
	h := new(handler)
	h.handler = subhandler
	h.auth = auth
	h.requiredGroups = requiredGroups
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("auth handler")
	if err := h.serveHTTP(resp, req); err != nil {
		status := http.StatusUnauthorized
		glog.V(2).Infof("auth failed => send %s", http.StatusText(status))
		resp.WriteHeader(status)
	} else {
		h.handler.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	name, value, err := header.ParseAuthorizationBearerHttpRequest(req)
	if err != nil {
		return err
	}
	token := header.CreateAuthorizationToken(name, value)
	glog.V(2).Infof("token: %s", token)
	user, err := h.auth(model.AuthToken(token), h.requiredGroups)
	if err != nil {
		glog.V(2).Infof("get user with token %s and group %v faild", token, h.requiredGroups)
		return err
	}
	glog.V(2).Infof("user %v is in group %v", *user, h.requiredGroups)
	return nil
}
