package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bborbe/auth/api"

	"github.com/bborbe/http/bearer"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type authClient struct {
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
	executeRequest             ExecuteRequest
	address                    string
	applicationName            api.ApplicationName
	applicationPassword        api.ApplicationPassword
}

type AuthClient interface {
	Auth(authToken api.AuthToken, requiredGroups []api.GroupName) (*api.UserName, error)
}

func New(executeRequest ExecuteRequest, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider, address string, applicationName api.ApplicationName, applicationPassword api.ApplicationPassword) *authClient {
	a := new(authClient)
	a.executeRequest = executeRequest
	a.httpRequestBuilderProvider = httpRequestBuilderProvider
	a.address = address
	a.applicationName = applicationName
	a.applicationPassword = applicationPassword
	return a
}

func (a *authClient) createBearer() string {
	return fmt.Sprintf("%s:%s", a.applicationName, a.applicationPassword)
}

func (a *authClient) Auth(authToken api.AuthToken, requiredGroups []api.GroupName) (*api.UserName, error) {
	request := api.LoginRequest{
		AuthToken:      authToken,
		RequiredGroups: requiredGroups,
	}
	target := fmt.Sprintf("http://%s/login", a.address)
	logger.Debugf("send request to %s", target)
	requestbuilder := a.httpRequestBuilderProvider.NewHttpRequestBuilder(target)
	requestbuilder.SetMethod("POST")
	requestbuilder.AddContentType("application/json")
	requestbuilder.AddHeader("Authorization", bearer.CreateBearerHeader(string(a.applicationName), string(a.applicationPassword)))
	content, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	logger.Debugf("auth request message: %s", string(content))
	requestbuilder.SetBody(bytes.NewBuffer(content))
	req, err := requestbuilder.Build()
	if err != nil {
		return nil, err
	}
	resp, err := a.executeRequest(req)
	if err != nil {
		logger.Debugf("auth request failed: %v", err)
		return nil, err
	}
	logger.Debugf("auth response status: %s", resp.Status)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request not success. status: %d", resp.Status)
	}
	responseContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logger.Debugf("response %s", string(responseContent))
	var response api.LoginResponse
	err = json.Unmarshal(responseContent, &response)
	if err != nil {
		return nil, err
	}
	return response.UserName, nil
}
