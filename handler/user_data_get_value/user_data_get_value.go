package user_data_get_value

import (
	"encoding/json"
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
)

var logger = log.DefaultLogger

type GetUserDataValue func(userName api.UserName, key string) (string, error)

type handler struct {
	getUserDataValue GetUserDataValue
}

func New(
	getUserDataValue GetUserDataValue,
) *handler {
	h := new(handler)
	h.getUserDataValue = getUserDataValue
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	logger.Debugf("getUserDataValue")
	if err := h.serveHTTP(resp, req); err != nil {
		logger.Debugf("Marshal json failed: %v", err)
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(resp, req)
	}
}

func (h *handler) serveHTTP(resp http.ResponseWriter, req *http.Request) error {
	logger.Debugf("getUserDataValue")
	var request api.GetUserDataValueRequest
	var response api.GetUserDataValueResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(response)
}
