package user

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/directory/token_user_directory"
	"github.com/bborbe/auth/directory/user_data_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
	"github.com/bborbe/auth/directory/user_token_directory"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type service struct {
	userTokenDirectory user_token_directory.UserTokenDirectory
	userGroupDirectory user_group_directory.UserGroupDirectory
	tokenUserDirectory token_user_directory.TokenUserDirectory
	userDataDirectory  user_data_directory.UserDataDirectory
}

type Service interface {
	DeleteUser(userName api.UserName) error
	DeleteUserWithToken(authToken api.AuthToken) error
	CreateUserWithToken(userName api.UserName, authToken api.AuthToken) error
	AddTokenToUserWithToken(newToken api.AuthToken, userToken api.AuthToken) error
	RemoveTokenFromUserWithToken(newToken api.AuthToken, userToken api.AuthToken) error
	VerifyTokenHasGroups(authToken api.AuthToken, requiredGroupNames []api.GroupName) (*api.UserName, error)
}

func New(
	userTokenDirectory user_token_directory.UserTokenDirectory,
	userGroupDirectory user_group_directory.UserGroupDirectory,
	tokenUserDirectory token_user_directory.TokenUserDirectory,
	userDataDirectory user_data_directory.UserDataDirectory,
) *service {
	s := new(service)
	s.userTokenDirectory = userTokenDirectory
	s.userGroupDirectory = userGroupDirectory
	s.tokenUserDirectory = tokenUserDirectory
	s.userDataDirectory = userDataDirectory
	return s
}

func (s *service) DeleteUserWithToken(authToken api.AuthToken) error {
	logger.Debugf("delete user with token %v", authToken)
	userName, err := s.tokenUserDirectory.FindUserByAuthToken(authToken)
	if err != nil {
		logger.Debugf("find user with token %v failed", authToken)
		return err
	}
	return s.DeleteUser(*userName)
}

func (s *service) DeleteUser(userName api.UserName) error {
	logger.Debugf("delete user %v", userName)
	tokens, err := s.userTokenDirectory.Get(userName)
	if err != nil {
		logger.Debugf("find tokens for user %v failed", userName)
		return err
	}
	for _, token := range *tokens {
		if err = s.tokenUserDirectory.Remove(token); err != nil {
			logger.Debugf("remove token %v failed", token)
		}
	}
	if err = s.userDataDirectory.Delete(userName); err != nil {
		logger.Debugf("remove user data %v failed", userName)
		return err
	}
	if err = s.userTokenDirectory.Delete(userName); err != nil {
		logger.Debugf("remove user %v failed", userName)
		return err
	}
	logger.Debugf("delete user %v successful", userName)
	return nil
}

func (h *service) CreateUserWithToken(userName api.UserName, authToken api.AuthToken) error {
	logger.Debugf("add token user %v with token %v", userName, authToken)
	if err := h.assertTokenNotUsed(authToken); err != nil {
		logger.Debugf("token %v already used", authToken)
		return err
	}
	if err := h.assertUserNameNotUser(userName); err != nil {
		logger.Debugf("userName %v already used", userName)
		return err
	}
	if err := h.tokenUserDirectory.Add(authToken, userName); err != nil {
		logger.Debugf("add user %v to token %v failed", userName, authToken)
		return err
	}
	if err := h.userTokenDirectory.Add(userName, authToken); err != nil {
		logger.Debugf("add token %v to user %v failed", authToken, userName)
		return err
	}
	logger.Debugf("add token %v to user %v successful", authToken, userName)
	return nil
}

func (h *service) assertTokenNotUsed(authToken api.AuthToken) error {
	logger.Debugf("assert token %s not used", authToken)
	exists, err := h.tokenUserDirectory.Exists(authToken)
	if err != nil {
		logger.Debugf("exists token failed: %v", err)
		return err
	}
	if exists {
		return fmt.Errorf("create user failed, token %s already used", authToken)
	}
	logger.Debugf("token not used")
	return nil
}

func (h *service) assertUserNameNotUser(userName api.UserName) error {
	logger.Debugf("assert user %s not existing", userName)
	exists, err := h.userTokenDirectory.Exists(userName)
	if err != nil {
		logger.Debugf("exists user failed: %v", err)
		return err
	}
	if exists {
		return fmt.Errorf("create user failed, user %s already exists", userName)
	}
	logger.Debugf("user not existing")
	return nil
}

func (h *service) AddTokenToUserWithToken(newToken api.AuthToken, userToken api.AuthToken) error {
	logger.Debugf("add token %v to user with token %v", newToken, userToken)
	if err := h.assertTokenNotUsed(newToken); err != nil {
		return err
	}
	userName, err := h.tokenUserDirectory.FindUserByAuthToken(userToken)
	if err != nil {
		return err
	}
	logger.Debugf("add token %v to user %v", newToken, *userName)
	if err := h.tokenUserDirectory.Add(newToken, *userName); err != nil {
		logger.Debugf("add token %v to user %v failed: %v", newToken, *userName, err)
		return err
	}
	if err := h.userTokenDirectory.Add(*userName, newToken); err != nil {
		logger.Debugf("add user %v to token %v failed: %v", *userName, newToken, err)
		return err
	}
	logger.Debugf("token added successful")
	return nil
}

func (h *service) RemoveTokenFromUserWithToken(newToken api.AuthToken, userToken api.AuthToken) error {
	logger.Debugf("remove token %v from user with token %v", newToken, userToken)
	userName, err := h.tokenUserDirectory.FindUserByAuthToken(userToken)
	if err != nil {
		return err
	}
	logger.Debugf("remove token %v from user %v", newToken, *userName)
	if err := h.tokenUserDirectory.Remove(newToken); err != nil {
		logger.Debugf("remove token %v failed: %v", newToken, err)
		return err
	}
	if err := h.userTokenDirectory.Remove(*userName, newToken); err != nil {
		logger.Debugf("remove token %v from user %v failed: %v", newToken, *userName, err)
		return err
	}
	logger.Debugf("token removed successful")
	return nil
}

func (s *service) VerifyTokenHasGroups(authToken api.AuthToken, requiredGroupNames []api.GroupName) (*api.UserName, error) {
	logger.Debugf("verify token %v has groups %v", authToken, requiredGroupNames)
	userName, err := s.tokenUserDirectory.FindUserByAuthToken(authToken)
	if err != nil {
		return userName, err
	}
	logger.Debugf("verify user %v has groups %v", *userName, requiredGroupNames)
	for _, groupName := range requiredGroupNames {
		containsGroup, err := s.userGroupDirectory.Contains(*userName, groupName)
		if err != nil {
			logger.Debugf("contains failed: %v", err)
			return userName, err
		}
		if !containsGroup {
			return userName, fmt.Errorf("user %v not in group %v", *userName, groupName)
		}
	}
	logger.Debugf("token %v has all required groups", authToken)
	return userName, nil
}
