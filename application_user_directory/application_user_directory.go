package application_user_directory

import "github.com/bborbe/auth/api"

type applicationUserDirectory struct {
}

type ApplicationUserDirectory interface {
	Add(applicationName api.ApplicationName, userName api.UserName) error
	Remove(applicationName api.ApplicationName, userName api.UserName) error
	Contains(applicationName api.ApplicationName, userName api.UserName) (bool, error)
}

func New() *applicationUserDirectory {
	return new(applicationUserDirectory)
}

func (a *applicationUserDirectory) Add(applicationName api.ApplicationName, userName api.UserName) error {
	return nil
}

func (a *applicationUserDirectory) Remove(applicationName api.ApplicationName, userName api.UserName) error {
	return nil
}

func (a *applicationUserDirectory) Contains(applicationName api.ApplicationName, userName api.UserName) (bool, error) {
	return true, nil
}
