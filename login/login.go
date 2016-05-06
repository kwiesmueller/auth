package login

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/http/bearer"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type CheckApplication func(api.ApplicationName, api.ApplicationPassword) error
type FindUserByAuthToken func(authToken api.AuthToken) (*api.UserName, error)
type IsUserNotFound func(err error) bool
type ApplicationContainsUser func(applicationName api.ApplicationName, userName api.UserName) (bool, error)
type ApplicationContainsGroupWithUser func(applicationName api.ApplicationName, groupName api.GroupName, userName api.UserName) (bool, error)

type handler struct {
	checkApplication                 CheckApplication
	findUserByAuthToken              FindUserByAuthToken
	isUserNotFound                   IsUserNotFound
	applicationContainsUser          ApplicationContainsUser
	applicationContainsGroupWithUser ApplicationContainsGroupWithUser
}

func New(checkApplication CheckApplication, findUserByAuthToken FindUserByAuthToken, isUserNotFound IsUserNotFound, applicationContainsUser ApplicationContainsUser, applicationContainsGroupWithUser ApplicationContainsGroupWithUser) *handler {
	h := new(handler)
	h.checkApplication = checkApplication
	h.findUserByAuthToken = findUserByAuthToken
	h.isUserNotFound = isUserNotFound
	h.applicationContainsUser = applicationContainsUser
	h.applicationContainsGroupWithUser = applicationContainsGroupWithUser
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
	applicationName, applicationPassword, err := parseApplication(req)
	if err != nil {
		return err
	}
	if err = h.checkApplication(applicationName, applicationPassword); err != nil {
		return err
	}
	var request api.LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	response, err := h.login(applicationName, request.AuthToken, request.RequiredGroups)
	if err != nil {
		if h.isUserNotFound(err) {
			logger.Infof("user not found: %s", request.AuthToken)
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return err
	}
	return json.NewEncoder(resp).Encode(response)
}

func parseApplication(req *http.Request) (api.ApplicationName, api.ApplicationPassword, error) {
	name, password, err := bearer.ParseBearerHttpRequest(req)
	if err != nil {
		return "", "", err
	}
	return api.ApplicationName(name), api.ApplicationPassword(password), nil
}

func (h *handler) login(applicationName api.ApplicationName, authToken api.AuthToken, requiredGroupNames []api.GroupName) (*api.LoginResponse, error) {
	logger.Debugf("login")
	userName, err := h.findUserByAuthToken(authToken)
	if err != nil {
		return nil, err
	}
	userFound, err := h.applicationContainsUser(applicationName, *userName)
	if err != nil {
		return nil, err
	}
	if !userFound {
		return nil, fmt.Errorf("user %s not found in application %s", userName, applicationName)
	}
	for _, groupName := range requiredGroupNames {
		containsGroup, err := h.applicationContainsGroupWithUser(applicationName, groupName, *userName)
		if err != nil {
			return nil, err
		}
		if !containsGroup {
			return nil, fmt.Errorf("user %s not in group %s in application %s", userName, groupName, applicationName)
		}
	}
	return &api.LoginResponse{
		User: userName,
	}, nil
}
