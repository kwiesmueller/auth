package user_group

import (
	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/directory/group_user_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type service struct {
	userGroupDirectory user_group_directory.UserGroupDirectory
	groupUserDirectory group_user_directory.GroupUserDirectory
}

type Service interface {
	AddUserToGroup(userName api.UserName, groupName api.GroupName) error
	RemoveUserFromGroup(userName api.UserName, groupName api.GroupName) error
}

func New(userGroupDirectory user_group_directory.UserGroupDirectory, groupUserDirectory group_user_directory.GroupUserDirectory) *service {
	s := new(service)
	s.userGroupDirectory = userGroupDirectory
	s.groupUserDirectory = groupUserDirectory
	return s
}

func (s *service) AddUserToGroup(userName api.UserName, groupName api.GroupName) error {
	logger.Debugf("add user %v to group %v", userName, groupName)
	if err := s.userGroupDirectory.Add(userName, groupName); err != nil {
		logger.Debugf("add user %v to group %v failed: %v", userName, groupName, err)
		return err
	}
	if err := s.groupUserDirectory.Add(groupName, userName); err != nil {
		logger.Debugf("add user %v to group %v failed: %v", userName, groupName, err)
		return err
	}
	logger.Debugf("added user %v to group %v successful", userName, groupName)
	return nil
}

func (s *service) RemoveUserFromGroup(userName api.UserName, groupName api.GroupName) error {
	logger.Debugf("remove user %v from group %v", userName, groupName)
	if err := s.userGroupDirectory.Remove(userName, groupName); err != nil {
		logger.Debugf("remove user %v from group %v failed: %v", userName, groupName, err)
		return err
	}
	if err := s.groupUserDirectory.Remove(groupName, userName); err != nil {
		logger.Debugf("remove user %v from group %v failed: %v", userName, groupName, err)
		return err
	}
	logger.Debugf("removed user %v from group %v successful", userName, groupName)
	return nil
}
