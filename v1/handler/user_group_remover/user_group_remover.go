package user_group_remover

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type RemoveUserFromGroup func(userName model.UserName, groupName model.GroupName) error

type handler struct {
	removeUserFromGroup RemoveUserFromGroup
}

func New(removeTokenToUserWithToken RemoveUserFromGroup) *handler {
	h := new(handler)
	h.removeUserFromGroup = removeTokenToUserWithToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("remove user from group")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("remove user from group failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("remove user from group success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.RemoveUserFromGroupRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.removeUserFromGroup(request.UserName, request.GroupName)
	if err != nil {
		return err
	}
	return nil
}
