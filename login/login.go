package login

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/user_directory"
	"github.com/bborbe/http/bearer"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type CheckApplication func(api.ApplicationName, api.ApplicationPassword) error

type handler struct {
	userDirectory    user_directory.UserDirectory
	checkApplication CheckApplication
}

func New(userDirectory user_directory.UserDirectory, checkApplication CheckApplication) *handler {
	h := new(handler)
	h.userDirectory = userDirectory
	h.checkApplication = checkApplication
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
	response, err := h.login(applicationName, &request)
	if err != nil {
		if h.userDirectory.IsUserNotFound(err) {
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

func (h *handler) login(applicationName api.ApplicationName, request *api.LoginRequest) (*api.LoginResponse, error) {
	logger.Debugf("login")
	user, err := h.userDirectory.FindUserByAuthToken(applicationName, request.AuthToken)
	if err != nil {
		return nil, err
	}
	groups, err := findGroupForUser(*user)
	if err != nil {
		return nil, err
	}
	return &api.LoginResponse{
		User:   user,
		Groups: groups,
	}, nil
}

func findGroupForUser(user api.UserName) (*[]api.Group, error) {
	if user == api.UserName("bborbe") {
		return &[]api.Group{api.Group("storage/admin")}, nil
	}
	return nil, nil
}
