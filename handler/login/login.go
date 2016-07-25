package login

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type VerifyTokenHasGroups func(authToken api.AuthToken, requiredGroupNames []api.GroupName) (*api.UserName, error)

type handler struct {
	verifyTokenHasGroups VerifyTokenHasGroups
}

func New(
	verifyTokenHasGroups VerifyTokenHasGroups,
) *handler {
	h := new(handler)
	h.verifyTokenHasGroups = verifyTokenHasGroups
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("login")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	logger.Debugf("login")
	var request api.LoginRequest
	var response api.LoginResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	userName, err := h.verifyTokenHasGroups(request.AuthToken, request.RequiredGroups)
	if err != nil {
		if userName == nil {
			logger.Infof("user not found: %s", request.AuthToken)
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return err
	}
	response.UserName = userName
	return json.NewEncoder(resp).Encode(response)
}
