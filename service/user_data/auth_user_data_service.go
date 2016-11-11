package user_data

import (
	"github.com/bborbe/auth/directory/username_data_directory"
	"github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

type userDataService struct {
	userDataDirectory username_data_directory.UsernameDataDirectory
}

func New(userDataDirectory username_data_directory.UsernameDataDirectory) *userDataService {
	s := new(userDataService)
	s.userDataDirectory = userDataDirectory
	return s
}

func (s *userDataService) Set(userName model.UserName, data map[string]string) error {
	glog.V(4).Infof("Set")
	return s.userDataDirectory.Set(userName, data)
}

func (s *userDataService) SetValue(userName model.UserName, key string, value string) error {
	glog.V(4).Infof("SetValue")
	return s.userDataDirectory.SetValue(userName, key, value)
}

func (s *userDataService) Get(userName model.UserName) (map[string]string, error) {
	glog.V(4).Infof("Get")
	return s.userDataDirectory.Get(userName)
}

func (s *userDataService) GetValue(userName model.UserName, key string) (string, error) {
	glog.V(4).Infof("GetValue")
	return s.userDataDirectory.GetValue(userName, key)
}

func (s *userDataService) Delete(userName model.UserName) error {
	glog.V(4).Infof("Delete")
	return s.userDataDirectory.Delete(userName)
}

func (s *userDataService) DeleteValue(userName model.UserName, key string) error {
	glog.V(4).Infof("DeleteValue")
	return s.userDataDirectory.DeleteValue(userName, key)
}
