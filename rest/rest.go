package rest

import (
	"fmt"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/http/header"
)

type callRest func(url string, method string, request interface{}, response interface{}, header http.Header) error

type Rest interface {
	Call(path string, method string, request interface{}, response interface{}) error
}

type rest struct {
	callRest            callRest
	url                 model.Url
	applicationName     model.ApplicationName
	applicationPassword model.ApplicationPassword
}

func New(
	callRest callRest,
	url model.Url,
	applicationName model.ApplicationName,
	applicationPassword model.ApplicationPassword,
) *rest {
	r := new(rest)
	r.callRest = callRest
	r.url = url
	r.applicationName = applicationName
	r.applicationPassword = applicationPassword
	return r
}

func (r *rest) Call(path string, method string, request interface{}, response interface{}) error {
	h := make(http.Header)
	h.Add("Authorization", header.CreateAuthorizationBearerHeader(r.applicationName.String(), r.applicationPassword.String()))
	return r.callRest(fmt.Sprintf("%s%s", r.url, path), method, request, response, h)
}
