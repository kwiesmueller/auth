package list

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/user_finder"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type handler struct {
	userFinder user_finder.UserFinder
}

func New(userFinder user_finder.UserFinder) *handler {
	h := new(handler)
	h.userFinder = userFinder
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
	var request api.Request
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	response, err := h.login(&request)
	if err != nil {
		if h.userFinder.IsUserNotFound(err) {
			logger.Infof("user not found: %s", request.AuthToken)
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return err
	}
	return json.NewEncoder(resp).Encode(response)
}

func (h *handler) login(request *api.Request) (*api.Response, error) {
	logger.Debugf("login")
	applicationId, err := findApplication(request.ApplicationName, request.ApplicationPassword)
	if err != nil {
		return nil, err
	}
	user, err := h.userFinder.FindUserByAuthToken(*applicationId, request.AuthToken)
	if err != nil {
		return nil, err
	}
	groups, err := findGroupForUser(*user)
	if err != nil {
		return nil, err
	}
	return &api.Response{
		User:   user,
		Groups: groups,
	}, nil
}

func findApplication(applicationName api.ApplicationName, applicationPassword api.ApplicationPassword) (*api.ApplicationId, error) {
	if applicationName == "name" && applicationPassword == "pw" {
		id := api.ApplicationId("adsf")
		return &id, nil
	}
	return nil, fmt.Errorf("application not found")
}

func findGroupForUser(user api.User) (*[]api.Group, error) {
	if user == api.User("bborbe") {
		return &[]api.Group{api.Group("storage/admin")}, nil
	}
	return nil, nil
}
