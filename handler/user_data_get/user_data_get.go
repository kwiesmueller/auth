package user_data_get

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
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
	var request api.GetUserDataRequest
	var response api.GetUserDataResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(response)
}
