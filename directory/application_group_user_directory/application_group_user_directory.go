package application_group_user_directory

import "github.com/bborbe/auth/api"

type applicationGroupUserDirectory struct {
}

type ApplicationGroupUserDirectory interface {
	Add(applicationName api.ApplicationName, groupName api.GroupName, userName api.UserName) error
	Remove(applicationName api.ApplicationName, groupName api.GroupName, userName api.UserName) error
	Contains(applicationName api.ApplicationName, groupName api.GroupName, userName api.UserName) (bool, error)
}

func New() *applicationGroupUserDirectory {
	return new(applicationGroupUserDirectory)
}

func (a *applicationGroupUserDirectory) Add(applicationName api.ApplicationName, groupName api.GroupName, userName api.UserName) error {
	return nil
}

func (a *applicationGroupUserDirectory) Remove(applicationName api.ApplicationName, groupName api.GroupName, userName api.UserName) error {
	return nil
}

func (a *applicationGroupUserDirectory) Contains(applicationName api.ApplicationName, groupName api.GroupName, userName api.UserName) (bool, error) {
	return true, nil
}
