package user_data

import (
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
