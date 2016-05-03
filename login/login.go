package list

import (
	"encoding/json"
	"net/http"

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
	user, err := h.userFinder.FindUserByAuthToken(request.AuthToken)
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

func findGroupForUser(user api.User) (*[]api.Group, error) {
	if user == api.User("bborbe") {
		return &[]api.Group{api.Group("storage/admin")}, nil
	}
	return nil, nil
}
