package user_group

import (
	"github.com/bborbe/auth/model"
)

type callRest func(path string, method string, request interface{}, response interface{}) error

type userGroupService struct {
	callRest callRest
}

func New(
	callRest callRest,
) *userGroupService {
	s := new(userGroupService)
	s.callRest = callRest
	return s
}

func (s *userGroupService) AddUserToGroup(userName model.UserName, groupName model.GroupName) error {
	panic("not implemented")
}

func (s *userGroupService) RemoveUserFromGroup(userName model.UserName, groupName model.GroupName) error {
	panic("not implemented")
}
