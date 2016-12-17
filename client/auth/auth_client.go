package auth

import (
	"net/http"
	"net/url"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	"github.com/golang/glog"
)

type callRest func(path string, values url.Values, method string, request interface{}, response interface{}) error

type authClient struct {
	callRest callRest
}

func New(
	callRest callRest,
) *authClient {
	u := new(authClient)
	u.callRest = callRest
	return u
}

func (u *authClient) HasGroups(authToken model.AuthToken, requiredGroups []model.GroupName) (bool, error) {
	glog.V(4).Infof("check user %v has groups %v", authToken, requiredGroups)
	userName, err := u.VerifyTokenHasGroups(authToken, requiredGroups)
	if err != nil {
		glog.V(2).Infof("check user %v has groups %v failed: %v", authToken, requiredGroups, err)
		return false, err
	}
	return userName != nil && len(*userName) > 0, nil
}

func (u *authClient) VerifyTokenHasGroups(authToken model.AuthToken, requiredGroupNames []model.GroupName) (*model.UserName, error) {
	glog.V(4).Infof("verify user with token %v has groups %v", authToken, requiredGroupNames)
	request := v1.LoginRequest{
		AuthToken:      authToken,
		RequiredGroups: requiredGroupNames,
	}
	var response v1.LoginResponse
	if err := u.callRest("/api/1.0/login", nil, http.MethodPost, &request, &response); err != nil {
		glog.V(2).Infof("verify user with token %v has groups %v failed: %v", authToken, requiredGroupNames, err)
		return nil, err
	}
	return response.UserName, nil
}
