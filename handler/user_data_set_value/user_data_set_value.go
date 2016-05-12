package user_data_set_value

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type SetUserDataValue func(userName api.UserName, key string, value string) error

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
	logger.Debugf("setUserDataValue")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	logger.Debugf("setUserDataValue")
	var request api.SetUserDataValueRequest
	var response api.SetUserDataValueResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(response)
}