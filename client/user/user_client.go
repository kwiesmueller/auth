package user

import (
	"net/http"

	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	"github.com/golang/glog"
)

type callRest func(path string, method string, request interface{}, response interface{}) error

type userService struct {
	callRest callRest
}

func New(
	callRest callRest,
) *userService {
	u := new(userService)
	u.callRest = callRest
	return u
}

func (u *userService) ListTokenOfUser(username model.UserName) ([]model.AuthToken, error) {
	glog.V(4).Infof("list tokens of user %v", username)
	var response []model.AuthToken
	if err := u.callRest(fmt.Sprintf("/api/1.0/token?username=%v", username), http.MethodGet, nil, &response); err != nil {
		glog.V(2).Infof("list tokens of user %v failed: %v", username, err)
		return nil, err
	}
	return response, nil
}

func (u *userService) HasGroups(authToken model.AuthToken, requiredGroups []model.GroupName) (bool, error) {
	glog.V(4).Infof("check user %v has groups %v", authToken, requiredGroups)
	userName, err := u.VerifyTokenHasGroups(authToken, requiredGroups)
	if err != nil {
		glog.V(2).Infof("check user %v has groups %v failed: %v", authToken, requiredGroups, err)
		return false, err
	}
	return userName != nil && len(*userName) > 0, nil
}

func (u *userService) VerifyTokenHasGroups(authToken model.AuthToken, requiredGroupNames []model.GroupName) (*model.UserName, error) {
	glog.V(4).Infof("verify user with token %v has groups %v", authToken, requiredGroupNames)
	request := v1.LoginRequest{
		AuthToken:      authToken,
		RequiredGroups: requiredGroupNames,
	}
	var response v1.LoginResponse
	if err := u.callRest("/api/1.0/login", http.MethodPost, &request, &response); err != nil {
		glog.V(2).Infof("verify user with token %v has groups %v failed: %v", authToken, requiredGroupNames, err)
		return nil, err
	}
	return response.UserName, nil
}

func (u *userService) List() ([]model.UserName, error) {
	glog.V(2).Infof("list usernames")
	var response []model.UserName
	if err := u.callRest("/api/1.0/user", http.MethodGet, nil, &response); err != nil {
		glog.V(2).Infof("list usernames failed: %v", err)
		return nil, err
	}
	return response, nil
}

func (u *userService) CreateUserWithToken(userName model.UserName, authToken model.AuthToken) error {
	glog.V(4).Infof("create user %s with token %s", userName, authToken)
	request := v1.RegisterRequest{
		AuthToken: model.AuthToken(authToken),
		UserName:  model.UserName(userName),
	}
	var response v1.RegisterResponse
	if err := u.callRest(fmt.Sprintf("/api/1.0/user"), "POST", &request, &response); err != nil {
		glog.V(2).Infof("create user %s failed: %v", userName, err)
		return err
	}
	glog.V(4).Infof("create user %s successful", userName)
	return nil
}

func (h *userService) AddTokenToUserWithToken(token model.AuthToken, authToken model.AuthToken) error {
	glog.V(4).Infof("add token %s to user with token %s", token, authToken)
	if authToken == token {
		return fmt.Errorf("token equals authToken")
	}
	request := v1.AddTokenRequest{
		AuthToken: model.AuthToken(authToken),
		Token:     model.AuthToken(token),
	}
	var response v1.AddTokenResponse
	if err := h.callRest("/api/1.0/token", "POST", &request, &response); err != nil {
		glog.V(2).Infof("add token failed: %v", err)
		return err
	}
	glog.V(4).Infof("add token successful")
	return nil
}

func (u *userService) RemoveTokenFromUserWithToken(token model.AuthToken, authToken model.AuthToken) error {
	glog.V(4).Infof("remove token %s to user with token %s", token, authToken)
	if authToken == token {
		return fmt.Errorf("token equals authToken")
	}
	request := v1.AddTokenRequest{
		AuthToken: model.AuthToken(authToken),
		Token:     model.AuthToken(token),
	}
	var response v1.AddTokenResponse
	if err := u.callRest("/api/1.0/token", "DELETE", &request, &response); err != nil {
		glog.V(2).Infof("remove token failed: %v", err)
		return err
	}
	glog.V(4).Infof("remove token successful")
	return nil
}

func (u *userService) DeleteUser(username model.UserName) error {
	glog.V(2).Infof("delete user %s", username)
	if err := u.callRest(fmt.Sprintf("/api/1.0/user/%s", username), "DELETE", nil, nil); err != nil {
		glog.V(2).Infof("delete user %s failed: %v", username, err)
		return err
	}
	glog.V(2).Infof("delete user %s successful", username)
	return nil
}

func (u *userService) DeleteUserWithToken(authToken model.AuthToken) error {
	glog.V(4).Infof("delete user with token %v", authToken)
	if err := u.callRest(fmt.Sprintf("/api/1.0/token/%v", authToken), "DELETE", nil, nil); err != nil {
		glog.V(2).Infof("delete user with token %s failed: %v", authToken, err)
		return err
	}
	glog.V(4).Infof("delete user with token %s successful", authToken)
	return nil
}

func (s *userService) AddTokenToUser(token model.AuthToken, username model.UserName) error {
	panic("not implemented")
}

func (s *userService) RemoveTokenFromUser(token model.AuthToken, username model.UserName) error {
	panic("not implemented")
}
