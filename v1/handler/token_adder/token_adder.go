package token_adder

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type AddTokenToUserWithToken func(newToken model.AuthToken, userToken model.AuthToken) error

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
	var request v1.AddTokenRequest
	var response v1.AddTokenResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	err := h.action(&request, &response)
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(&response)
}

func (h *handler) action(request *v1.AddTokenRequest, response *v1.AddTokenResponse) error {
	return h.addTokenToUserWithToken(request.Token, request.AuthToken)
}
