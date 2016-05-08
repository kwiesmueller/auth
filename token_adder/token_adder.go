package token_adder

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type ExistsToken func(authToken api.AuthToken) (bool, error)
type FindUserByAuthToken func(authToken api.AuthToken) (*api.UserName, error)
type UserAddToken func(userName api.UserName, authToken api.AuthToken) error
type TokenAddUser func(authToken api.AuthToken, userName api.UserName) error

type handler struct {
	userAddToken        UserAddToken
	tokenAddUser        TokenAddUser
	existsToken         ExistsToken
	findUserByAuthToken FindUserByAuthToken
}

func New(
	userAddToken UserAddToken,
	tokenAddUser TokenAddUser,
	existsToken ExistsToken,
	findUserByAuthToken FindUserByAuthToken,
) *handler {
	h := new(handler)
	h.userAddToken = userAddToken
	h.tokenAddUser = tokenAddUser
	h.existsToken = existsToken
	h.findUserByAuthToken = findUserByAuthToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("add token")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("add token failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("add token success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.AddTokenRequest
	var response api.AddTokenResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *api.AddTokenRequest, response *api.AddTokenResponse) error {
	if err := h.assertTokenNotUsed(request.Token); err != nil {
		return err
	}
	userName, err := h.findUserByAuthToken(request.AuthToken)
	if err != nil {
		return err
	}
	if err := h.tokenAddUser(request.Token, *userName); err != nil {
		return err
	}
	if err := h.userAddToken(*userName, request.Token); err != nil {
		return err
	}
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
