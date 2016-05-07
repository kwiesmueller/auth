package user_creator

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
	logger.Debugf("create application")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.RegisterRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	if err := h.assertTokenNotUsed(request.AuthToken); err != nil {
		return err
	}
	if err := h.assertUserNameNotUser(request.UserName); err != nil {
		return err
	}
	if err := h.tokenAddUser(request.AuthToken, request.UserName); err != nil {
		return err
	}
	if err := h.userAddToken(request.UserName, request.AuthToken); err != nil {
		return err
	}
	return nil
}

func (h *handler) assertTokenNotUsed(authToken api.AuthToken) error {
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
