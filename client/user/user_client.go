package user

import (
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
)

type callRest func(path string, method string, request interface{}, response interface{}) error

type userService struct {
	callRest callRest
}

func New(
	callRest callRest,
) *userService {
	s := new(userService)
	s.callRest = callRest
	return s
}

func (s *userService) HasGroups(authToken model.AuthToken, requiredGroups []model.GroupName) (bool, error) {
	panic("not implemented")
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
	request := v1.LoginRequest{
		AuthToken:      authToken,
		RequiredGroups: requiredGroupNames,
	}
	var response v1.LoginResponse
	if err := s.callRest("/api/1.0/login", http.MethodPost, &request, &response); err != nil {
		return nil, err
	}
	return response.UserName, nil
}

func (s *userService) List() ([]model.UserName, error) {
	panic("not implemented")
}
