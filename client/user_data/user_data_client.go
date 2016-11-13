package user_data

import (
	"fmt"

	"github.com/bborbe/auth/model"
)

type callRest func(path string, method string, request interface{}, response interface{}) error

type userDataService struct {
	callRest callRest
}

func New(
	callRest callRest,
) *userDataService {
	s := new(userDataService)
	s.callRest = callRest
	return s
}

func (s *userDataService) Set(userName model.UserName, data map[string]string) error {
	return fmt.Errorf("not implemented")
}

func (s *userDataService) SetValue(userName model.UserName, key string, value string) error {
	return fmt.Errorf("not implemented")
}

func (s *userDataService) Get(userName model.UserName) (map[string]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *userDataService) GetValue(userName model.UserName, key string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *userDataService) Delete(userName model.UserName) error {
	return fmt.Errorf("not implemented")
}

func (s *userDataService) DeleteValue(userName model.UserName, key string) error {
	return fmt.Errorf("not implemented")
}
