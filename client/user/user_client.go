package user

import (
	"github.com/bborbe/auth/model"
)

type userService struct {
}

func New() *userService {
	s := new(userService)
	return s
}

func (s *userService) DeleteUserWithToken(authToken model.AuthToken) error {
	panic("not implemented")
}

func (s *userService) DeleteUser(userName model.UserName) error {
	panic("not implemented")
}

func (h *userService) CreateUserWithToken(userName model.UserName, authToken model.AuthToken) error {
	panic("not implemented")
}

func (h *userService) assertTokenNotUsed(authToken model.AuthToken) error {
	panic("not implemented")
}

func (h *userService) AddTokenToUserWithToken(newToken model.AuthToken, userToken model.AuthToken) error {
	panic("not implemented")
}

func (h *userService) RemoveTokenFromUserWithToken(newToken model.AuthToken, userToken model.AuthToken) error {
	panic("not implemented")
}

func (s *userService) VerifyTokenHasGroups(authToken model.AuthToken, requiredGroupNames []model.GroupName) (*model.UserName, error) {
	panic("not implemented")
}

func (s *userService) List() ([]model.UserName, error) {
	panic("not implemented")
}
