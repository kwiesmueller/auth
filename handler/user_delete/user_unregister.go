package user_delete

import (
	"net/http"

	"fmt"
	"strings"

	"github.com/bborbe/auth/api"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type DeleteUser func(username api.UserName) error

type handler struct {
	deleteUser DeleteUser
}

func New(
	deleteUser DeleteUser,
) *handler {
	h := new(handler)
	h.deleteUser = deleteUser
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("delete user")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("delete user failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		logger.Debugf("delete user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	parts := strings.Split(req.RequestURI, "/")
	if len(parts) == 0 {
		logger.Debugf("auth token missing")
		return fmt.Errorf("invalid request uri: %s", req.RequestURI)
	}
	username := api.UserName(parts[len(parts)-1])
	if err := h.deleteUser(username); err != nil {
		return err
	}
	resp.WriteHeader(http.StatusOK)
	return nil
}
