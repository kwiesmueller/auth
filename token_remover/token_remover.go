package token_remover

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type FindUserByAuthToken func(authToken api.AuthToken) (*api.UserName, error)
type UserRemoveToken func(userName api.UserName, authToken api.AuthToken) error
type TokenRemoveUser func(authToken api.AuthToken, userName api.UserName) error

type handler struct {
	userRemoveToken     UserRemoveToken
	tokenRemoveUser     TokenRemoveUser
	findUserByAuthToken FindUserByAuthToken
}

func New(
	userRemoveToken UserRemoveToken,
	tokenRemoveUser TokenRemoveUser,
	findUserByAuthToken FindUserByAuthToken,
) *handler {
	h := new(handler)
	h.userRemoveToken = userRemoveToken
	h.tokenRemoveUser = tokenRemoveUser
	h.findUserByAuthToken = findUserByAuthToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("remove token")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("remove token failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.RemoveTokenRequest
	logger.Debugf("decode json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	userName, err := h.findUserByAuthToken(request.AuthToken)
	if err != nil {
		return err
	}
	if err := h.tokenRemoveUser(request.Token, *userName); err != nil {
		return err
	}
	if err := h.userRemoveToken(*userName, request.Token); err != nil {
		return err
	}
	return nil
}
