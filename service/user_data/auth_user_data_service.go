package user_data

import (
	"github.com/bborbe/auth/directory/user_data_directory"
	"github.com/bborbe/auth/model"
)

type userDataService struct {
	userDataDirectory user_data_directory.UserDataDirectory
}

func New(userDataDirectory user_data_directory.UserDataDirectory) *userDataService {
	s := new(userDataService)
	s.userDataDirectory = userDataDirectory
	return s
}

func (s *userDataService) Set(userName model.UserName, data map[string]string) error {
	return s.userDataDirectory.Set(userName, data)
}

func (s *userDataService) SetValue(userName model.UserName, key string, value string) error {
	return s.userDataDirectory.SetValue(userName, key, value)
}

func (s *userDataService) Get(userName model.UserName) (map[string]string, error) {
	return s.userDataDirectory.Get(userName)
}

func (s *userDataService) GetValue(userName model.UserName, key string) (string, error) {
	return s.userDataDirectory.GetValue(userName, key)
}

func (s *userDataService) Delete(userName model.UserName) error {
	return s.userDataDirectory.Delete(userName)
}

func (s *userDataService) DeleteValue(userName model.UserName, key string) error {
	return s.userDataDirectory.DeleteValue(userName, key)
}
