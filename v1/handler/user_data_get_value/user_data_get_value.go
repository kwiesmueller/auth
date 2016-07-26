package user_data_get_value

import (
	"encoding/json"
	"net/http"

	"fmt"
	"regexp"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type GetUserDataValue func(userName model.UserName, key string) (string, error)

type handler struct {
	getUserDataValue GetUserDataValue
}

func New(
	getUserDataValue GetUserDataValue,
) *handler {
	h := new(handler)
	h.getUserDataValue = getUserDataValue
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("getUserDataValue")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	logger.Debugf("getUserData")
	path := req.URL.Path
	logger.Debugf("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data/(.*)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 3 {
		return fmt.Errorf("find user failed")
	}
	data, err := h.getUserDataValue(model.UserName(matches[1]), matches[2])
	if err != nil {
		return err
	}
	response := v1.GetUserDataValueResponse(data)
	return json.NewEncoder(resp).Encode(&response)
}
