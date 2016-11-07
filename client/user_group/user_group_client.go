package user_group

import (
	"github.com/bborbe/auth/model"
)

type userGroupService struct {
}

func New() *userGroupService {
	s := new(userGroupService)
	return s
}

func (s *userGroupService) AddUserToGroup(userName model.UserName, groupName model.GroupName) error {
	panic("not implemented")
}

func (s *userGroupService) RemoveUserFromGroup(userName model.UserName, groupName model.GroupName) error {
	panic("not implemented")
}
