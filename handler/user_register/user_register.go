package user_register

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type CreateUserWithToken func(userName api.UserName, authToken api.AuthToken) error

type handler struct {
	createUserWithToken CreateUserWithToken
}

func New(addTokenToUser CreateUserWithToken) *handler {
	h := new(handler)
	h.createUserWithToken = addTokenToUser
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("register user")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("register user failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("register user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.RegisterRequest
	var response api.RegisterResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *api.RegisterRequest, response *api.RegisterResponse) error {
	return h.createUserWithToken(request.UserName, request.AuthToken)
}
