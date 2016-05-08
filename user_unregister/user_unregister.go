package user_unregister

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type FindUserByAuthToken func(authToken api.AuthToken) (*api.UserName, error)
type GetTokensForUser func(userName api.UserName) (*[]api.AuthToken, error)
type RemoveToken func(authToken api.AuthToken) error
type RemoveUser func(userName api.UserName) error

type handler struct {
	findUserByAuthToken FindUserByAuthToken
	getTokensForUser    GetTokensForUser
	removeToken         RemoveToken
	removeUser          RemoveUser
}

func New(
	findUserByAuthToken FindUserByAuthToken,
	getTokensForUser GetTokensForUser,
	removeToken RemoveToken,
	removeUser RemoveUser,
) *handler {
	h := new(handler)
	h.findUserByAuthToken = findUserByAuthToken
	h.getTokensForUser = getTokensForUser
	h.removeToken = removeToken
	h.removeUser = removeUser
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("user create")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("create user failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.UnRegisterRequest
	logger.Debugf("decode json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	logger.Debugf("unregister user with token %v", request.AuthToken)
	userName, err := h.findUserByAuthToken(request.AuthToken)
	if err != nil {
		logger.Debugf("find user with token %v failed", request.AuthToken)
		return err
	}
	tokens, err := h.getTokensForUser(*userName)
	if err != nil {
		logger.Debugf("find tokens for user %v failed", *userName)
		return err
	}
	for _, token := range *tokens {
		if err = h.removeToken(token); err != nil {
			logger.Debugf("remove token %v failed", token)
		}
	}
	if err = h.removeUser(*userName); err != nil {
		logger.Debugf("remove user %v failed", *userName)
		return err
	}
	logger.Debugf("unregister user %v successful", *userName)
	return nil
}
