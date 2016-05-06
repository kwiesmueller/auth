package application_group_directory

import "github.com/bborbe/auth/api"

type applicationGroupDirectory struct {
}

type ApplicationGroupDirectory interface {
	Add(applicationName api.ApplicationName, groupName api.GroupName) error
	Remove(applicationName api.ApplicationName, groupName api.GroupName) error
	Contains(applicationName api.ApplicationName, groupName api.GroupName) (bool, error)
}

func New() *applicationGroupDirectory {
	return new(applicationGroupDirectory)
}

func (a *applicationGroupDirectory) Add(applicationName api.ApplicationName, groupName api.GroupName) error {
	return nil
}

func (a *applicationGroupDirectory) Remove(applicationName api.ApplicationName, groupName api.GroupName) error {
	return nil
}

func (a *applicationGroupDirectory) Contains(applicationName api.ApplicationName, groupName api.GroupName) (bool, error) {
	return true, nil
}
