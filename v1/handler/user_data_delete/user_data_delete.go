package user_data_delete

import (
	"net/http"

	"fmt"
	"regexp"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type DeleteUserData func(userName model.UserName) error

type handler struct {
	deleteUserData DeleteUserData
}

func New(
	deleteUserData DeleteUserData,
) *handler {
	h := new(handler)
	h.deleteUserData = deleteUserData
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("deleteUserData")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	glog.V(2).Infof("deleteUserData")
	path := req.URL.Path
	glog.V(2).Infof("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 2 {
		return fmt.Errorf("find user failed")
	}
	return h.deleteUserData(model.UserName(matches[1]))
}
