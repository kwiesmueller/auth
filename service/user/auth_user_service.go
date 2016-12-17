package user

import (
	"fmt"

	"github.com/bborbe/auth/directory/token_username_directory"
	"github.com/bborbe/auth/directory/username_data_directory"
	"github.com/bborbe/auth/directory/username_groupname_directory"
	"github.com/bborbe/auth/directory/username_token_directory"
	"github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

type userService struct {
	usernameTokenDirectory     username_token_directory.UsernameTokenDirectory
	usernameGroupnameDirectory username_groupname_directory.UsernameGroupnameDirectory
	tokenUsernameDirectory     token_username_directory.TokenUsernameDirectory
	usernameDataDirectory      username_data_directory.UsernameDataDirectory
}

func New(
	usernameTokenDirectory username_token_directory.UsernameTokenDirectory,
	usernameGroupnameDirectory username_groupname_directory.UsernameGroupnameDirectory,
	tokenUsernameDirectory token_username_directory.TokenUsernameDirectory,
	usernameDataDirectory username_data_directory.UsernameDataDirectory,
) *userService {
	s := new(userService)
	s.usernameTokenDirectory = usernameTokenDirectory
	s.usernameGroupnameDirectory = usernameGroupnameDirectory
	s.tokenUsernameDirectory = tokenUsernameDirectory
	s.usernameDataDirectory = usernameDataDirectory
	return s
}

func (u *userService) DeleteUserWithToken(authToken model.AuthToken) error {
	glog.V(4).Infof("delete user with token %v", authToken)
	username, err := u.tokenUsernameDirectory.FindUserByAuthToken(authToken)
	if err != nil {
		glog.V(2).Infof("find user with token %v failed", authToken)
		return err
	}
	return u.DeleteUser(*username)
}

func (u *userService) DeleteUser(username model.UserName) error {
	glog.V(4).Infof("delete user %v", username)
	tokens, err := u.usernameTokenDirectory.Get(username)
	if err != nil {
		glog.V(2).Infof("find tokens for user %v failed", username)
		return err
	}
	for _, token := range tokens {
		if err = u.tokenUsernameDirectory.Remove(token); err != nil {
			glog.V(2).Infof("remove token %v failed", token)
		}
	}
	if err = u.usernameDataDirectory.Delete(username); err != nil {
		glog.V(2).Infof("remove user data %v failed", username)
		return err
	}
	if err = u.usernameTokenDirectory.Delete(username); err != nil {
		glog.V(2).Infof("remove user %v failed", username)
		return err
	}
	glog.V(4).Infof("delete user %v successful", username)
	return nil
}

func (u *userService) CreateUserWithToken(username model.UserName, authToken model.AuthToken) error {
	glog.V(4).Infof("add token %v to user %v", authToken, username)
	if len(username) == 0 {
		return fmt.Errorf("username empty")
	}
	if len(authToken) == 0 {
		return fmt.Errorf("token empty")
	}
	if err := u.assertTokenNotUsed(authToken); err != nil {
		glog.V(2).Infof("token %v already used", authToken)
		return err
	}
	if err := u.assertUserNameNotUsed(username); err != nil {
		glog.V(2).Infof("username %v already used", username)
		return err
	}
	if err := u.tokenUsernameDirectory.Add(authToken, username); err != nil {
		glog.V(2).Infof("add user %v to token %v failed", username, authToken)
		return err
	}
	if err := u.usernameTokenDirectory.Add(username, authToken); err != nil {
		glog.V(2).Infof("add token %v to user %v failed", authToken, username)
		return err
	}
	glog.V(4).Infof("add token %v to user %v successful", authToken, username)
	return nil
}

func (u *userService) assertTokenNotUsed(authToken model.AuthToken) error {
	glog.V(4).Infof("assert token %s not used", authToken)
	exists, err := u.tokenUsernameDirectory.Exists(authToken)
	if err != nil {
		glog.V(2).Infof("exists token failed: %v", err)
		return err
	}
	if exists {
		return fmt.Errorf("token %s already used", authToken)
	}
	glog.V(4).Infof("token not used")
	return nil
}

func (u *userService) assertUserNameNotUsed(username model.UserName) error {
	glog.V(4).Infof("assert user %s not existing", username)
	exists, err := u.usernameTokenDirectory.Exists(username)
	if err != nil {
		glog.V(2).Infof("exists user failed: %v", err)
		return err
	}
	if exists {
		return fmt.Errorf("create user failed, user %s already exists", username)
	}
	glog.V(4).Infof("user not existing")
	return nil
}

func (u *userService) AddTokenToUserWithToken(newToken model.AuthToken, userToken model.AuthToken) error {
	glog.V(4).Infof("add token %v to user with token %v", newToken, userToken)
	if err := u.assertTokenNotUsed(newToken); err != nil {
		glog.V(2).Infof("token %v already used, can't add token", err)
		return err
	}
	return u.AddTokenToUserWithTokenForce(newToken, userToken)
}

func (u *userService) AddTokenToUserWithTokenForce(newToken model.AuthToken, userToken model.AuthToken) error {
	glog.V(4).Infof("add token %v to user with token %v", newToken, userToken)
	username, err := u.tokenUsernameDirectory.FindUserByAuthToken(userToken)
	if err != nil {
		glog.V(2).Infof("find user with token %v failed: %v", userToken, err)
		return err
	}
	return u.AddTokenToUser(newToken, *username)
}

func (u *userService) RemoveTokenFromUserWithToken(newToken model.AuthToken, userToken model.AuthToken) error {
	glog.V(4).Infof("remove token %v from user with token %v", newToken, userToken)
	username, err := u.tokenUsernameDirectory.FindUserByAuthToken(userToken)
	if err != nil {
		glog.V(2).Infof("find user with token %v failed: %v", userToken, err)
		return err
	}
	return u.RemoveTokenFromUser(newToken, *username)
}

func (u *userService) List() ([]model.UserName, error) {
	glog.V(4).Infof("list users")
	result, err := u.usernameTokenDirectory.List()
	if err != nil {
		glog.V(4).Infof("list users failed: %v", err)
		return nil, err
	}
	glog.V(4).Infof("found %d users", len(result))
	return result, nil
}

func (u *userService) ListTokenOfUser(username model.UserName) ([]model.AuthToken, error) {
	glog.V(4).Infof("list tokens of user %v", username)
	result, err := u.usernameTokenDirectory.Get(username)
	if err != nil {
		glog.V(2).Infof("list tokens of user %v failed: %v", username, err)
		return nil, err
	}
	glog.V(4).Infof("found %d tokens for user %v", len(result), username)
	return result, nil
}

func (u *userService) AddTokenToUser(newToken model.AuthToken, username model.UserName) error {
	glog.V(2).Infof("add token %v to user %v", newToken, username)
	if err := u.tokenUsernameDirectory.Add(newToken, username); err != nil {
		glog.V(2).Infof("add token %v to user %v failed: %v", newToken, username, err)
		return err
	}
	if err := u.usernameTokenDirectory.Add(username, newToken); err != nil {
		glog.V(2).Infof("add user %v to token %v failed: %v", username, newToken, err)
		return err
	}
	glog.V(4).Infof("token added successful")
	return nil
}

func (u *userService) RemoveTokenFromUser(newToken model.AuthToken, username model.UserName) error {
	glog.V(4).Infof("remove token %v from user %v", newToken, username)
	if err := u.tokenUsernameDirectory.Remove(newToken); err != nil {
		glog.V(2).Infof("remove token %v failed: %v", newToken, err)
		return err
	}
	if err := u.usernameTokenDirectory.Remove(username, newToken); err != nil {
		glog.V(2).Infof("remove token %v from user %v failed: %v", newToken, username, err)
		return err
	}
	glog.V(4).Infof("token removed successful")
	return nil
}
