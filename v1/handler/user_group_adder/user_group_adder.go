package user_group_adder

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type AddUserToGroup func(userName model.UserName, groupName model.GroupName) error

type handler struct {
	addUserToGroup AddUserToGroup
}

func New(addUserToGroup AddUserToGroup) *handler {
	h := new(handler)
	h.addUserToGroup = addUserToGroup
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("add user to group")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("add user to group failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(2).Infof("add user to group success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request v1.AddUserToGroupRequest
	var response v1.AddUserToGroupResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *v1.AddUserToGroupRequest, response *v1.AddUserToGroupResponse) error {
	return h.addUserToGroup(request.UserName, request.GroupName)
}
