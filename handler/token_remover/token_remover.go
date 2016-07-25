package token_remover

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type RemoveTokenToUserWithToken func(newToken api.AuthToken, userToken api.AuthToken) error

type handler struct {
	removeTokenToUserWithToken RemoveTokenToUserWithToken
}

func New(removeTokenToUserWithToken RemoveTokenToUserWithToken) *handler {
	h := new(handler)
	h.removeTokenToUserWithToken = removeTokenToUserWithToken
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
	return h.removeTokenToUserWithToken(request.Token, request.AuthToken)
}
