package login

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type verifyTokenHasGroups func(authToken model.AuthToken, requiredGroupNames []model.GroupName) (*model.UserName, error)

type handler struct {
	verifyTokenHasGroups verifyTokenHasGroups
}

func New(
	verifyTokenHasGroups verifyTokenHasGroups,
) *handler {
	h := new(handler)
	h.verifyTokenHasGroups = verifyTokenHasGroups
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("login")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("login failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
		return
	}
	glog.V(3).Infof("login success")
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.LoginRequest
	var response v1.LoginResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		glog.V(3).Infof("decode json failed: %v", err)
		return err
	}
	glog.V(4).Infof("verify user with token %v has groups %v", request.AuthToken, request.RequiredGroups)
	userName, err := h.verifyTokenHasGroups(request.AuthToken, request.RequiredGroups)
	if err != nil {
		glog.V(3).Infof("verify token has group failed: %v", err)
		if userName == nil {
			glog.V(1).Infof("user not found: %s", request.AuthToken)
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return err
	}
	glog.V(4).Infof("user with token %v has groups %v", request.AuthToken, request.RequiredGroups)
	response.UserName = userName
	return json.NewEncoder(resp).Encode(response)
}
