package user_unregister

import (
	"net/http"

	"fmt"
	"strings"

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
	logger.Debugf("unregister user")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("unregister user failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("unregister user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	parts := strings.Split(req.RequestURI, "/")
	if len(parts) == 0 {
		logger.Debugf("auth token missing")
		return fmt.Errorf("invalid request uri: %s", req.RequestURI)
	}
	authToken := api.AuthToken(parts[len(parts)-1])
	logger.Debugf("unregister user with token %v", authToken)
	userName, err := h.findUserByAuthToken(authToken)
	if err != nil {
		logger.Debugf("find user with token %v failed", authToken)
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
