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
	userService    service.UserService
	requiredGroups []model.GroupName
}

func New(
	userService service.UserService,
	requiredGroups ...model.GroupName,
) *handler {
	h := new(handler)
	h.userService = userService
	h.requiredGroups = requiredGroups
	return h
}

func (h *handler) Check(req *http.Request) (bool, error) {
	name, value, err := header.ParseAuthorizationBearerHttpRequest(req)
	if err != nil {
		glog.V(2).Infof("parse authorization header failed: %v", err)
		return false, err
	}
	return h.userService.HasGroups(
		model.AuthToken(header.CreateAuthorizationToken(name, value)),
		h.requiredGroups,
	)
}
