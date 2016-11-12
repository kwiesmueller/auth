package auth

import (
	"fmt"

	"github.com/bborbe/auth/directory/token_username_directory"
	"github.com/bborbe/auth/directory/username_groupname_directory"
	"github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

type userService struct {
	usernameGroupnameDirectory username_groupname_directory.UsernameGroupnameDirectory
	tokenUsernameDirectory     token_username_directory.TokenUsernameDirectory
}

func New(
	usernameGroupnameDirectory username_groupname_directory.UsernameGroupnameDirectory,
	tokenUsernameDirectory token_username_directory.TokenUsernameDirectory,
) *userService {
	s := new(userService)
	s.usernameGroupnameDirectory = usernameGroupnameDirectory
	s.tokenUsernameDirectory = tokenUsernameDirectory
	return s
}

func (u *userService) VerifyTokenHasGroups(authToken model.AuthToken, requiredGroupNames []model.GroupName) (*model.UserName, error) {
	glog.V(4).Infof("verify token %v has groups %v", authToken, requiredGroupNames)
	username, err := u.tokenUsernameDirectory.FindUserByAuthToken(authToken)
	if err != nil {
		glog.V(2).Infof("find user by token failed: %v", err)
		return nil, err
	}
	glog.V(4).Infof("verify user %v has groups %v", *username, requiredGroupNames)
	for _, groupName := range requiredGroupNames {
		containsGroup, err := u.usernameGroupnameDirectory.Contains(*username, groupName)
		if err != nil {
			glog.V(2).Infof("contains failed: %v", err)
			return username, err
		}
		if !containsGroup {
			return username, fmt.Errorf("user %v not in group %v", *username, groupName)
		}
	}
	glog.V(4).Infof("token %v has all required groups", authToken)
	return username, nil
}

func (u *userService) HasGroups(authToken model.AuthToken, requiredGroups []model.GroupName) (bool, error) {
	glog.V(4).Infof("check user with token %v has groups %v", authToken, requiredGroups)
	username, err := u.VerifyTokenHasGroups(authToken, requiredGroups)
	if err != nil {
		glog.V(2).Infof("check user with token %v has groups %v failed: %v", authToken, requiredGroups, err)
		return false, err
	}
	return username != nil && len(*username) > 0, nil
}
