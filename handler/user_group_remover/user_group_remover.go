package user_group_remover

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type RemoveUserFromGroup func(userName api.UserName, groupName api.GroupName) error

type handler struct {
	removeUserFromGroup RemoveUserFromGroup
}

func New(removeTokenToUserWithToken RemoveUserFromGroup) *handler {
	h := new(handler)
	h.removeUserFromGroup = removeTokenToUserWithToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("remove user from group")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("remove user from group failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("remove user from group success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.RemoveUserFromGroupRequest
	var response api.RemoveUserFromGroupResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *api.RemoveUserFromGroupRequest, response *api.RemoveUserFromGroupResponse) error {
	return h.removeUserFromGroup(request.UserName, request.GroupName)
}
