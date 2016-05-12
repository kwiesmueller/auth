package user_data_set

import (
	"encoding/json"
	"net/http"

	"fmt"
	"regexp"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type SetUserData func(userName api.UserName, data map[string]string) error

type handler struct {
	setUserData SetUserData
}

func New(
	setUserData SetUserData,
) *handler {
	h := new(handler)
	h.setUserData = setUserData
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("setUserData")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	logger.Debugf("setUserData")
	var request api.SetUserDataRequest
	path := req.URL.Path
	logger.Debugf("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 2 {
		return fmt.Errorf("find user failed")
	}
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	return h.setUserData(api.UserName(matches[1]), request)
}
