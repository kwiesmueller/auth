package user

import (
	"fmt"

	"github.com/bborbe/auth/directory/token_user_directory"
	"github.com/bborbe/auth/directory/user_data_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
	"github.com/bborbe/auth/directory/user_token_directory"
	"github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

type userService struct {
	userTokenDirectory user_token_directory.UserTokenDirectory
	userGroupDirectory user_group_directory.UserGroupDirectory
	tokenUserDirectory token_user_directory.TokenUserDirectory
	userDataDirectory  user_data_directory.UserDataDirectory
}

func New(
	userTokenDirectory user_token_directory.UserTokenDirectory,
	userGroupDirectory user_group_directory.UserGroupDirectory,
	tokenUserDirectory token_user_directory.TokenUserDirectory,
	userDataDirectory user_data_directory.UserDataDirectory,
) *userService {
	s := new(userService)
	s.userTokenDirectory = userTokenDirectory
	s.userGroupDirectory = userGroupDirectory
	s.tokenUserDirectory = tokenUserDirectory
	s.userDataDirectory = userDataDirectory
	return s
}

func (s *userService) DeleteUserWithToken(authToken model.AuthToken) error {
	glog.V(2).Infof("delete user with token %v", authToken)
	userName, err := s.tokenUserDirectory.FindUserByAuthToken(authToken)
	if err != nil {
		glog.V(2).Infof("find user with token %v failed", authToken)
		return err
	}
	return s.DeleteUser(*userName)
}

func (s *userService) DeleteUser(userName model.UserName) error {
	glog.V(2).Infof("delete user %v", userName)
	tokens, err := s.userTokenDirectory.Get(userName)
	if err != nil {
		glog.V(2).Infof("find tokens for user %v failed", userName)
		return err
	}
	for _, token := range tokens {
		if err = s.tokenUserDirectory.Remove(token); err != nil {
			glog.V(2).Infof("remove token %v failed", token)
		}
	}
	if err = s.userDataDirectory.Delete(userName); err != nil {
		glog.V(2).Infof("remove user data %v failed", userName)
		return err
	}
	if err = s.userTokenDirectory.Delete(userName); err != nil {
		glog.V(2).Infof("remove user %v failed", userName)
		return err
	}
	glog.V(2).Infof("delete user %v successful", userName)
	return nil
}

func (h *userService) CreateUserWithToken(userName model.UserName, authToken model.AuthToken) error {
	glog.V(2).Infof("add token user %v with token %v", userName, authToken)
	if err := h.assertTokenNotUsed(authToken); err != nil {
		glog.V(2).Infof("token %v already used", authToken)
		return err
	}
	if err := h.assertUserNameNotUser(userName); err != nil {
		glog.V(2).Infof("userName %v already used", userName)
		return err
	}
	if err := h.tokenUserDirectory.Add(authToken, userName); err != nil {
		glog.V(2).Infof("add user %v to token %v failed", userName, authToken)
		return err
	}
	if err := h.userTokenDirectory.Add(userName, authToken); err != nil {
		glog.V(2).Infof("add token %v to user %v failed", authToken, userName)
		return err
	}
	glog.V(2).Infof("add token %v to user %v successful", authToken, userName)
	return nil
}

func (h *userService) assertTokenNotUsed(authToken model.AuthToken) error {
	glog.V(4).Infof("assert token %s not used", authToken)
	exists, err := h.tokenUserDirectory.Exists(authToken)
	if err != nil {
		glog.V(2).Infof("exists token failed: %v", err)
		return err
	}
	if exists {
		return fmt.Errorf("token %s already used", authToken)
	}
	glog.V(2).Infof("token not used")
	return nil
}

func (h *userService) assertUserNameNotUser(userName model.UserName) error {
	glog.V(4).Infof("assert user %s not existing", userName)
	exists, err := h.userTokenDirectory.Exists(userName)
	if err != nil {
		glog.V(2).Infof("exists user failed: %v", err)
		return err
	}
	if exists {
		return fmt.Errorf("create user failed, user %s already exists", userName)
	}
	glog.V(2).Infof("user not existing")
	return nil
}

func (h *userService) AddTokenToUserWithToken(newToken model.AuthToken, userToken model.AuthToken) error {
	glog.V(4).Infof("add token %v to user with token %v", newToken, userToken)
	if err := h.assertTokenNotUsed(newToken); err != nil {
		glog.V(2).Infof("token %v already used, can't add token", err)
		return err
	}
	return h.AddTokenToUserWithTokenForce(newToken, userToken)
}

func (h *userService) AddTokenToUserWithTokenForce(newToken model.AuthToken, userToken model.AuthToken) error {
	glog.V(4).Infof("add token %v to user with token %v", newToken, userToken)
	userName, err := h.tokenUserDirectory.FindUserByAuthToken(userToken)
	if err != nil {
		glog.V(2).Infof("find user with token %v failed: %v", userToken, err)
		return err
	}
	glog.V(2).Infof("add token %v to user %v", newToken, *userName)
	if err := h.tokenUserDirectory.Add(newToken, *userName); err != nil {
		glog.V(2).Infof("add token %v to user %v failed: %v", newToken, *userName, err)
		return err
	}
	if err := h.userTokenDirectory.Add(*userName, newToken); err != nil {
		glog.V(2).Infof("add user %v to token %v failed: %v", *userName, newToken, err)
		return err
	}
	glog.V(4).Infof("token added successful")
	return nil
}

func (h *userService) RemoveTokenFromUserWithToken(newToken model.AuthToken, userToken model.AuthToken) error {
	glog.V(4).Infof("remove token %v from user with token %v", newToken, userToken)
	userName, err := h.tokenUserDirectory.FindUserByAuthToken(userToken)
	if err != nil {
		glog.V(2).Infof("find user with token %v failed: %v", userToken, err)
		return err
	}
	glog.V(4).Infof("remove token %v from user %v", newToken, *userName)
	if err := h.tokenUserDirectory.Remove(newToken); err != nil {
		glog.V(2).Infof("remove token %v failed: %v", newToken, err)
		return err
	}
	if err := h.userTokenDirectory.Remove(*userName, newToken); err != nil {
		glog.V(2).Infof("remove token %v from user %v failed: %v", newToken, *userName, err)
		return err
	}
	glog.V(2).Infof("token removed successful")
	return nil
}

func (s *userService) VerifyTokenHasGroups(authToken model.AuthToken, requiredGroupNames []model.GroupName) (*model.UserName, error) {
	glog.V(4).Infof("verify token %v has groups %v", authToken, requiredGroupNames)
	userName, err := s.tokenUserDirectory.FindUserByAuthToken(authToken)
	if err != nil {
		glog.V(2).Infof("find user by token failed: %v", err)
		return nil, err
	}
	glog.V(4).Infof("verify user %v has groups %v", *userName, requiredGroupNames)
	for _, groupName := range requiredGroupNames {
		containsGroup, err := s.userGroupDirectory.Contains(*userName, groupName)
		if err != nil {
			glog.V(2).Infof("contains failed: %v", err)
			return userName, err
		}
		if !containsGroup {
			return userName, fmt.Errorf("user %v not in group %v", *userName, groupName)
		}
	}
	glog.V(4).Infof("token %v has all required groups", authToken)
	return userName, nil
}

func (s *userService) List() ([]model.UserName, error) {
	glog.V(4).Infof("list users")
	result, err := s.userTokenDirectory.List()
	if err != nil {
		glog.V(4).Infof("list users failed: %v", err)
		return nil, err
	}
	glog.V(4).Infof("found %d users", len(result))
	return result, nil
}

func (s *userService) HasGroups(authToken model.AuthToken, requiredGroups []model.GroupName) (bool, error) {
	glog.V(4).Infof("check user with token %v has groups %v", authToken, requiredGroups)
	userName, err := s.VerifyTokenHasGroups(authToken, requiredGroups)
	if err != nil {
		glog.V(2).Infof("check user with token %v has groups %v failed: %v", authToken, requiredGroups, err)
		return false, err
	}
	return userName != nil && len(*userName) > 0, nil
}

func (s *userService) ListTokenOfUser(username model.UserName) ([]model.AuthToken, error) {
	glog.V(4).Infof("list tokens of user %v", username)
	result, err := s.userTokenDirectory.Get(username)
	if err != nil {
		glog.V(2).Infof("list tokens of user %v failed: %v", username, err)
		return nil, err
	}
	glog.V(4).Infof("found %d tokens for user %v", len(result), username)
	return result, nil
}
