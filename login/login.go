package list

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

const (
	NOT_FOUND = fmt.Errorf("user not found")
)

var logger = log.DefaultLogger

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
	var request Request
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	response, err := login(&request)
	if err != nil {
		if err == NOT_FOUND {
			logger.Infof("user not found: %s/%s", request.ConnectorName, request.ConnectorUserIdentifier)
			resp.WriteHeader(http.StatusNotFound)
			return nil
		}
		return err
	}
	return json.NewEncoder(resp).Encode(response)
}

func login(request *Request) (*Reponse, error) {
	logger.Debugf("login")
	user, err := findUser(request.ConnectorName, request.ConnectorUserIdentifier)
	if err != nil {
		return nil, err
	}
	groups, err := findGroupForUser(user)
	if err != nil {
		return nil, err
	}
	return &Reponse{
		User:   user,
		Groups: groups,
	}
}

func findUser(connectorName string, userIdentifier string) (User, error) {
	logger.Debugf("find user with connector: %s and userId: %s", connectorName, userIdentifier)
	if connectorName == "hipchat" && userIdentifier == "130647" {
		return User("bborbe")
	}
	if connectorName == "telegram" && userIdentifier == "asda" {
		return User("bborbe")
	}
	return NOT_FOUND
}

func findGroupForUser(user User) ([]Group, error) {
	if user == User("bborbe") {
		return []Group{Group("storage/admin")}, nil
	}
	return nil, nil
}

type User string

type Group string

type Request struct {
	ApplicationName         string `json:"applicatonName"`
	ApplicationPassword     string `json:"applicatonPassword"`
	ConnectorName           string `json:"connectorName"`
	ConnectorUserIdentifier string `json:"connectorUserIdentifier"`
}

type Reponse struct {
	User   User     `json:"user"`
	Groups []string `json:"groups"`
}
