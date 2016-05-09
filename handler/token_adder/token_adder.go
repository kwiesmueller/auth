package token_adder

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type AddTokenToUserWithToken func(newToken api.AuthToken, userToken api.AuthToken) error

type handler struct {
	addTokenToUserWithToken AddTokenToUserWithToken
}

func New(addTokenToUserWithToken AddTokenToUserWithToken) *handler {
	h := new(handler)
	h.addTokenToUserWithToken = addTokenToUserWithToken
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
	return h.addTokenToUserWithToken(request.Token, request.AuthToken)
}
