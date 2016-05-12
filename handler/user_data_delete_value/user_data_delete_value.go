package user_data_delete_value

import (
	"encoding/json"
	"net/http"

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
	logger.Debugf("deleteUserDataValue")
	var request api.DeleteUserDataValueRequest
	var response api.DeleteUserDataValueResponse
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return err
	}
	return json.NewEncoder(resp).Encode(response)
}
