package user_list

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type ListUsers func() ([]model.UserName, error)

type handler struct {
	listUsers ListUsers
}

func New(removeTokenToUserWithToken ListUsers) *handler {
	h := new(handler)
	h.listUsers = removeTokenToUserWithToken
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
	var response v1.UserListResponse
	var err error
	var userNames []model.UserName
	if userNames, err = h.listUsers(); err != nil {
		return err
	}
	for _, userName := range userNames {
		response = append(response, v1.User{UserName: userName})
	}
	return json.NewEncoder(resp).Encode(&response)
}
