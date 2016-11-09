package tokens_by_username

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/bborbe/auth/model"
	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type listAuthTokenOfUser func(model.UserName) ([]model.AuthToken, error)

const parameter = "username"

type handler struct {
	listAuthTokenOfUser listAuthTokenOfUser
}

func New(listAuthTokenOfUser listAuthTokenOfUser) *handler {
	h := new(handler)
	h.listAuthTokenOfUser = listAuthTokenOfUser
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	glog.V(4).Infof("list tokens for user")
	if err := h.serveHTTP(resp, req); err != nil {
		glog.V(2).Infof("list tokens for user failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	} else {
		glog.V(4).Infof("list tokens for user success")
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	username := model.UserName(req.FormValue(parameter))
	glog.V(4).Infof("list tokens for user %v", username)
	if len(username) == 0 {
		glog.V(2).Infof("parameter %v missing", parameter)
		return fmt.Errorf("parameter %v missing", parameter)
	}
	var err error
	var result []model.AuthToken
	if result, err = h.listAuthTokenOfUser(username); err != nil {
		glog.V(2).Infof("list tokens for user %v: failed: %v", username, err)
		return err
	}
	glog.V(2).Infof("got %d tokens for user %v", len(result), username)
	return json.NewEncoder(resp).Encode(&result)
}
