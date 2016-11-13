package user_data_get

import (
	"net/http"

	"encoding/json"
	"fmt"
	"regexp"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type GetUserData func(userName model.UserName) (map[string]string, error)

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
	glog.V(2).Infof("getUserData")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("Marshal json failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	glog.V(2).Infof("getUserData")
	path := req.URL.Path
	glog.V(2).Infof("path: %s", path)
	re := regexp.MustCompile(`/user/([^/]*)/data`)
	matches := re.FindStringSubmatch(path)
	if len(matches) != 2 {
		return fmt.Errorf("find user failed")
	}
	data, err := h.getUserData(model.UserName(matches[1]))
	if err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(data)
}
