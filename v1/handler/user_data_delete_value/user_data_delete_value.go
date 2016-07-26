package user_data_delete_value

import (
	"net/http"

	"fmt"
	"regexp"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type DeleteUserDataValue func(userName model.UserName, key string) error

type handler struct {
	deleteUserDataValue DeleteUserDataValue
}

func New(
	deleteUserDataValue DeleteUserDataValue,
) *handler {
	h := new(handler)
	h.deleteUserDataValue = deleteUserDataValue
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("deleteUserDataValue")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	logger.Debugf("deleteUserDataValue")
	path := req.URL.Path
	logger.Debugf("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data/(.*)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 3 {
		return fmt.Errorf("find user failed")
	}
	return h.deleteUserDataValue(model.UserName(matches[1]), matches[2])
}
