package user_data_set

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

type SetUserData func(userName model.UserName, data map[string]string) error

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
	glog.V(2).Infof("setUserData")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	glog.V(2).Infof("setUserData")
	var request v1.SetUserDataRequest
	path := req.URL.Path
	glog.V(2).Infof("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 2 {
		return fmt.Errorf("find user failed")
	}
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	return h.setUserData(model.UserName(matches[1]), request)
}
