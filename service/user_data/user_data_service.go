package user_data

import (
	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/directory/user_data_directory"
)

type UserDataService interface {
	Set(userName api.UserName, data map[string]string) error
	SetValue(userName api.UserName, key string, value string) error
	Get(userName api.UserName) (map[string]string, error)
	GetValue(userName api.UserName, key string) (string, error)
	Delete(userName api.UserName) error
	DeleteValue(userName api.UserName, key string) error
}

type service struct {
	userDataDirectory user_data_directory.UserDataDirectory
}

func New(userDataDirectory user_data_directory.UserDataDirectory) *service {
	s := new(service)
	s.userDataDirectory = userDataDirectory
	return s
}

func (s *service) Set(userName api.UserName, data map[string]string) error {
	return s.userDataDirectory.Set(userName, data)
}

func (s *service) SetValue(userName api.UserName, key string, value string) error {
	return s.userDataDirectory.SetValue(userName, key, value)
}

func (s *service) Get(userName api.UserName) (map[string]string, error) {
	return s.userDataDirectory.Get(userName)
}

func (s *service) GetValue(userName api.UserName, key string) (string, error) {
	return s.userDataDirectory.GetValue(userName, key)
}

func (s *service) Delete(userName api.UserName) error {
	return s.userDataDirectory.Delete(userName)
}

func (s *service) DeleteValue(userName api.UserName, key string) error {
	return s.userDataDirectory.DeleteValue(userName, key)
}
