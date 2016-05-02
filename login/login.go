package list

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var (
	NOT_FOUND = fmt.Errorf("user not found")
	logger    = log.DefaultLogger
)

type handler struct {
}

func New() *handler {
	h := new(handler)
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("create")
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
	response, err := login(&request)
	if err != nil {
		if err == NOT_FOUND {
			logger.Infof("user not found: %s", request.AuthToken)
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return err
	}
	return json.NewEncoder(resp).Encode(response)
}

func login(request *api.Request) (*api.Response, error) {
	logger.Debugf("login")
	user, err := findUser(request.AuthToken)
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

func findUser(authToken api.AuthToken) (*api.User, error) {
	logger.Debugf("find user with auth token: %s", authToken)
	if authToken == api.AuthToken("hipchat:130647") {
		user := api.User("bborbe")
		return &user, nil
	}
	if authToken == api.AuthToken("telegram:abc") {
		user := api.User("bborbe")
		return &user, nil
	}
	return nil, NOT_FOUND
}

func findGroupForUser(user api.User) (*[]api.Group, error) {
	if user == api.User("bborbe") {
		return &[]api.Group{api.Group("storage/admin")}, nil
	}
	return nil, nil
}
