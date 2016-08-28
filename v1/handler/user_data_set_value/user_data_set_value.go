package user_data_set_value

import (
	"encoding/json"
	"net/http"

	"fmt"
	"regexp"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type SetUserDataValue func(userName model.UserName, key string, value string) error

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
	glog.V(2).Infof("setUserDataValue")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	glog.V(2).Infof("setUserDataValue")
	var request v1.SetUserDataValueRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	path := req.URL.Path
	glog.V(2).Infof("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data/(.*)`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 3 {
		return fmt.Errorf("find user failed")
	}
	return h.setUserDataValue(model.UserName(matches[1]), matches[2], string(request))
}
