package verify_group_handler

import (
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/http/header"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type Auth func(authToken api.AuthToken, requiredGroups []api.GroupName) (*api.UserName, error)

type handler struct {
	auth           Auth
	handler        http.Handler
	requiredGroups []api.GroupName
}

func New(subhandler http.Handler, auth Auth, requiredGroups ...api.GroupName) *handler {
	h := new(handler)
	h.handler = subhandler
	h.auth = auth
	h.requiredGroups = requiredGroups
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("auth handler")
	if err := h.serveHTTP(resp, req); err != nil {
		status := http.StatusUnauthorized
		logger.Debugf("auth failed => send %s", http.StatusText(status))
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
	logger.Debugf("token: %s", token)
	user, err := h.auth(api.AuthToken(token), h.requiredGroups)
	if err != nil {
		logger.Debugf("get user with token %s and group %v faild", token, h.requiredGroups)
		return err
	}
	logger.Debugf("user %v is in group %v", *user, h.requiredGroups)
	return nil
}
