package client

import (
	"fmt"
	"net/http"

	"github.com/bborbe/auth/client/application"
	"github.com/bborbe/auth/client/auth"
	"github.com/bborbe/auth/client/user"
	"github.com/bborbe/auth/client/user_data"
	"github.com/bborbe/auth/client/user_group"
	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/service"
	"github.com/bborbe/http/header"
	"github.com/bborbe/http/rest"
	"net/url"
)

type authClient struct {
	executeRequest      func(req *http.Request) (resp *http.Response, err error)
	url                 model.Url
	applicationName     model.ApplicationName
	applicationPassword model.ApplicationPassword
}

type Client interface {
	AuthService() service.AuthService
	ApplicationService() service.ApplicationService
	UserDataService() service.UserDataService
	UserGroupService() service.UserGroupService
	UserService() service.UserService
}

func New(
executeRequest func(req *http.Request) (resp *http.Response, err error),
url model.Url,
applicationName model.ApplicationName,
applicationPassword model.ApplicationPassword,
) *authClient {
	r := new(authClient)
	r.executeRequest = executeRequest
	r.url = url
	r.applicationName = applicationName
	r.applicationPassword = applicationPassword
	return r
}

func (r *authClient) call(path string, values url.Values, method string, request interface{}, response interface{}) error {
	h := make(http.Header)
	h.Add("Authorization", header.CreateAuthorizationBearerHeader(r.applicationName.String(), r.applicationPassword.String()))
	return rest.New(r.executeRequest).Call(fmt.Sprintf("%s%s", r.url, path), values, method, request, response, h)
}

func (r *authClient) ApplicationService() service.ApplicationService {
	return application.New(r.call)
}

func (r *authClient) UserDataService() service.UserDataService {
	return user_data.New(r.call)
}

func (r *authClient) UserGroupService() service.UserGroupService {
	return user_group.New(r.call)
}

func (r *authClient) UserService() service.UserService {
	return user.New(r.call)
}

func (r *authClient) AuthService() service.AuthService {
	return auth.New(r.call)
}
