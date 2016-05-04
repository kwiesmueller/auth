package application_directory

import "github.com/bborbe/auth/api"

type applicationDirectory struct {
}

type ApplicationDirectory interface {
	Check(api.ApplicationName, api.ApplicationPassword) error
}

func New() *applicationDirectory {
	return new(applicationDirectory)
}

func (a *applicationDirectory) Check(api.ApplicationName, api.ApplicationPassword) error {
	return nil
}
