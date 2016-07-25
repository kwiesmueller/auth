package user_data_get

import (
	"net/http"

	"encoding/json"
	"fmt"
	"regexp"

	"github.com/bborbe/auth/api"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type GetUserData func(userName api.UserName) (map[string]string, error)

type handler struct {
	getUserData GetUserData
}

func New(
	getUserData GetUserData,
) *handler {
	h := new(handler)
	h.getUserData = getUserData
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("getUserData")
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
	re := regexp.MustCompile(`/user/([^/]*)/data`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 2 {
		return fmt.Errorf("find user failed")
	}
	data, err := h.getUserData(api.UserName(matches[1]))
	if err != nil {
		return err
	}
	response := api.GetUserDataResponse(data)
	return json.NewEncoder(resp).Encode(&response)
}
