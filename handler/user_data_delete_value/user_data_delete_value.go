package user_data_delete_value

import (
	"net/http"

	"fmt"
	"regexp"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type DeleteUserDataValue func(userName api.UserName, key string) error

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
	logger.Debugf("deleteUserData")
	path := req.URL.Path
	logger.Debugf("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data/(.*)`)
	matches := re.FindStringSubmatch(path)
	fmt.Printf("%v", matches)
	if len(matches) != 3 {
		return fmt.Errorf("find user failed")
	}
	return h.deleteUserDataValue(api.UserName(matches[1]), matches[2])
}
