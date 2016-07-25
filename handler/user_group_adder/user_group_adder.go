package user_group_adder

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type AddUserToGroup func(userName api.UserName, groupName api.GroupName) error

type handler struct {
	addUserToGroup AddUserToGroup
}

func New(addUserToGroup AddUserToGroup) *handler {
	h := new(handler)
	h.addUserToGroup = addUserToGroup
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("add user to group")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("add user to group failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("add user to group success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.AddUserToGroupRequest
	var response api.AddUserToGroupResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *api.AddUserToGroupRequest, response *api.AddUserToGroupResponse) error {
	return h.addUserToGroup(request.UserName, request.GroupName)
}
