package user_data

import (
	"github.com/bborbe/auth/directory/user_data_directory"
	"github.com/bborbe/auth/model"
)

type UserDataService interface {
	Set(userName model.UserName, data map[string]string) error
	SetValue(userName model.UserName, key string, value string) error
	Get(userName model.UserName) (map[string]string, error)
	GetValue(userName model.UserName, key string) (string, error)
	Delete(userName model.UserName) error
	DeleteValue(userName model.UserName, key string) error
	List() ([]model.UserName, error)
}

type service struct {
	userDataDirectory user_data_directory.UserDataDirectory
}

func New(userDataDirectory user_data_directory.UserDataDirectory) *service {
	s := new(service)
	s.userDataDirectory = userDataDirectory
	return s
}

func (s *service) Set(userName model.UserName, data map[string]string) error {
	return s.userDataDirectory.Set(userName, data)
}

func (s *service) SetValue(userName model.UserName, key string, value string) error {
	return s.userDataDirectory.SetValue(userName, key, value)
}

func (s *service) Get(userName model.UserName) (map[string]string, error) {
	return s.userDataDirectory.Get(userName)
}

func (s *service) GetValue(userName model.UserName, key string) (string, error) {
	return s.userDataDirectory.GetValue(userName, key)
}

func (s *service) Delete(userName model.UserName) error {
	return s.userDataDirectory.Delete(userName)
}

func (s *service) DeleteValue(userName model.UserName, key string) error {
	return s.userDataDirectory.DeleteValue(userName, key)
}

func (s *service) List() ([]model.UserName, error) {
	return nil, nil
}
