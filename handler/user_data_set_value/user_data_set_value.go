package user_data_set_value

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

type SetUserDataValue func(userName api.UserName, key string, value string) error

type handler struct {
	setUserDataValue SetUserDataValue
}

func New(
	setUserDataValue SetUserDataValue,
) *handler {
	h := new(handler)
	h.setUserDataValue = setUserDataValue
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("setUserDataValue")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	logger.Debugf("setUserDataValue")
	var request api.SetUserDataValueRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	path := req.URL.Path
	logger.Debugf("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data/(.*)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 3 {
		return fmt.Errorf("find user failed")
	}
	return h.setUserDataValue(api.UserName(matches[1]), matches[2], string(request))
}
