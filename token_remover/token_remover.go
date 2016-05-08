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
type TokenRemove func(authToken api.AuthToken) error

type handler struct {
	userRemoveToken     UserRemoveToken
	tokenRemove         TokenRemove
	findUserByAuthToken FindUserByAuthToken
}

func New(
	userRemoveToken UserRemoveToken,
	tokenRemove TokenRemove,
	findUserByAuthToken FindUserByAuthToken,
) *handler {
	h := new(handler)
	h.userRemoveToken = userRemoveToken
	h.tokenRemove = tokenRemove
	h.findUserByAuthToken = findUserByAuthToken
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("remove token")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("remove token failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("remove token success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	var request api.RemoveTokenRequest
	var response api.RemoveTokenResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *api.RemoveTokenRequest, response *api.RemoveTokenResponse) error {
	userName, err := h.findUserByAuthToken(request.AuthToken)
	if err != nil {
		return err
	}
	if err := h.tokenRemove(request.Token); err != nil {
		return err
	}
	if err := h.userRemoveToken(*userName, request.Token); err != nil {
		return err
	}
	return nil
}
