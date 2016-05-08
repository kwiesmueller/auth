package user_register

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type UserAddToken func(userName api.UserName, authToken api.AuthToken) error
type TokenAddUser func(authToken api.AuthToken, userName api.UserName) error
type ExistsUser func(userName api.UserName) (bool, error)
type ExistsToken func(authToken api.AuthToken) (bool, error)

type handler struct {
	userAddToken UserAddToken
	tokenAddUser TokenAddUser
	existsUser   ExistsUser
	existsToken  ExistsToken
}

func New(
	userAddToken UserAddToken,
	tokenAddUser TokenAddUser,
	existsUser ExistsUser,
	existsToken ExistsToken,
) *handler {
	h := new(handler)
	h.userAddToken = userAddToken
	h.tokenAddUser = tokenAddUser
	h.existsUser = existsUser
	h.existsToken = existsToken
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
	logger.Debugf("register user %v with token %v", request.UserName, request.AuthToken)
	if err := h.assertTokenNotUsed(request.AuthToken); err != nil {
		logger.Debugf("token %v already used", request.AuthToken)
		return err
	}
	if err := h.assertUserNameNotUser(request.UserName); err != nil {
		logger.Debugf("userName %v already used", request.UserName)
		return err
	}
	if err := h.tokenAddUser(request.AuthToken, request.UserName); err != nil {
		logger.Debugf("add user %v to token %v failed", request.UserName, request.AuthToken)
		return err
	}
	if err := h.userAddToken(request.UserName, request.AuthToken); err != nil {
		logger.Debugf("add token %v to user %v failed", request.AuthToken, request.UserName)
		return err
	}
	logger.Debugf("register user %v successful", request.UserName)
	return nil
}

func (h *handler) assertTokenNotUsed(authToken api.AuthToken) error {
	logger.Debugf("assert token %s not used", authToken)
	exists, err := h.existsToken(authToken)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("create user failed, token %s already used", authToken)
	}
	logger.Debugf("token not used")
	return nil
}

func (h *handler) assertUserNameNotUser(userName api.UserName) error {
	logger.Debugf("assert user %s not existing", userName)
	exists, err := h.existsUser(userName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("create user failed, user %s already exists", userName)
	}
	logger.Debugf("user not existing")
	return nil
}
