package user_data

import (
	"github.com/bborbe/auth/model"
)

type userDataService struct {
}

func New() *userDataService {
	s := new(userDataService)
	return s
}

func (s *userDataService) Set(userName model.UserName, data map[string]string) error {
	panic("not implemented")
}

func (s *userDataService) SetValue(userName model.UserName, key string, value string) error {
	panic("not implemented")
}

func (s *userDataService) Get(userName model.UserName) (map[string]string, error) {
	panic("not implemented")
}

func (s *userDataService) GetValue(userName model.UserName, key string) (string, error) {
	panic("not implemented")
}

func (s *userDataService) Delete(userName model.UserName) error {
	panic("not implemented")
}

func (s *userDataService) DeleteValue(userName model.UserName, key string) error {
	panic("not implemented")
}
