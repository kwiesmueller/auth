package user_group

import (
	"github.com/bborbe/auth/directory/group_user_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
)

type service struct {
	userGroupDirectory user_group_directory.UserGroupDirectory
	groupUserDirectory group_user_directory.GroupUserDirectory
}
type Service interface {
}

func New(userGroupDirectory user_group_directory.UserGroupDirectory, groupUserDirectory group_user_directory.GroupUserDirectory) *service {
	s := new(service)
	s.userGroupDirectory = userGroupDirectory
	s.groupUserDirectory = groupUserDirectory
	return s
}
