package check

import (
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/service"
	"github.com/bborbe/http/header"
	"github.com/golang/glog"
)

type hasGroups func(authToken model.AuthToken, requiredGroups []model.GroupName) (bool, error)

type Check interface {
	Check(req *http.Request) (bool, error)
}

type handler struct {
	authService    service.AuthService
	requiredGroups []model.GroupName
}

func New(
	authService service.AuthService,
	requiredGroups ...model.GroupName,
) *handler {
	h := new(handler)
	h.authService = authService
	h.requiredGroups = requiredGroups
	return h
}

func (h *handler) Check(req *http.Request) (bool, error) {
	name, value, err := header.ParseAuthorizationBearerHttpRequest(req)
	if err != nil {
		glog.V(2).Infof("parse authorization header failed: %v", err)
		return false, err
	}
	return h.authService.HasGroups(
		model.AuthToken(header.CreateAuthorizationToken(name, value)),
		h.requiredGroups,
	)
}
