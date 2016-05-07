package user_creator

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type handler struct {
}

func New() *handler {
	h := new(handler)
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

	return nil
}

func (h *handler) assertTokenNotUsed(authToken api.AuthToken) error {
	return nil
}

func (h *handler) assertUserNameNotUser(userName api.UserName) error {
	return nil
}
