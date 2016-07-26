package user_list

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

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
	logger.Debugf("list user")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("list user failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("list user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var response []v1.User
	userNames, err := h.listUsers()
	if err != nil {
		return err
	}
	for _, userName := range userNames {
		response = append(response, v1.User{UserName: userName})
	}
	return json.NewEncoder(resp).Encode(&response)
}
