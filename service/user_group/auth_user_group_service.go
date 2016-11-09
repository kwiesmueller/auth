package user_group

import (
	"github.com/bborbe/auth/directory/group_user_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
	"github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

type userGroupService struct {
	userGroupDirectory user_group_directory.UserGroupDirectory
	groupUserDirectory group_user_directory.GroupUserDirectory
}

func New(userGroupDirectory user_group_directory.UserGroupDirectory, groupUserDirectory group_user_directory.GroupUserDirectory) *userGroupService {
	s := new(userGroupService)
	s.userGroupDirectory = userGroupDirectory
	s.groupUserDirectory = groupUserDirectory
	return s
}

func (s *userGroupService) AddUserToGroup(userName model.UserName, groupName model.GroupName) error {
	glog.V(4).Infof("add user %v to group %v", userName, groupName)
	if err := s.userGroupDirectory.Add(userName, groupName); err != nil {
		glog.V(2).Infof("add user %v to group %v failed: %v", userName, groupName, err)
		return err
	}
	if err := s.groupUserDirectory.Add(groupName, userName); err != nil {
		glog.V(2).Infof("add user %v to group %v failed: %v", userName, groupName, err)
		return err
	}
	glog.V(4).Infof("added user %v to group %v successful", userName, groupName)
	return nil
}

func (s *userGroupService) RemoveUserFromGroup(userName model.UserName, groupName model.GroupName) error {
	glog.V(4).Infof("remove user %v from group %v", userName, groupName)
	if err := s.userGroupDirectory.Remove(userName, groupName); err != nil {
		glog.V(2).Infof("remove user %v from group %v failed: %v", userName, groupName, err)
		return err
	}
	if err := s.groupUserDirectory.Remove(groupName, userName); err != nil {
		glog.V(2).Infof("remove user %v from group %v failed: %v", userName, groupName, err)
		return err
	}
	glog.V(4).Infof("removed user %v from group %v successful", userName, groupName)
	return nil
}

func (s *userGroupService) ListGroupNamesForUsername(username model.UserName) ([]model.GroupName, error) {
	glog.V(4).Infof("get groupNames for user %v", username)
	result, err := s.userGroupDirectory.Get(username)
	if err != nil {
		glog.V(2).Infof("get groupNames for user %v failed: %v", username, err)
		return nil, err
	}
	glog.V(4).Infof("got %d groupNames for user %v", len(result), username)
	return result, nil
}
