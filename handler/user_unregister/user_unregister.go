package user_unregister

import (
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/api"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type DeleteUserWithToken func(authToken api.AuthToken) error

type handler struct {
	deleteUserWithToken DeleteUserWithToken
}

func New(
	deleteUserWithToken DeleteUserWithToken,
) *handler {
	h := new(handler)
	h.deleteUserWithToken = deleteUserWithToken
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
	if err := h.deleteUserWithToken(authToken); err != nil {
		return err
	}
	resp.WriteHeader(http.StatusOK)
	return nil
}
