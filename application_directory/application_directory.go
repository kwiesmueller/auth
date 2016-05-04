package application_directory

import (
	"github.com/bborbe/auth/api"
	"github.com/timelinelabs/romulus/Godeps/_workspace/src/github.com/albertrdixon/gearbox/logger"
)

const AUTH_APPLICATION_NAME = api.ApplicationName("auth")

type applicationDirectory struct {
}

type ApplicationDirectory interface {
	Check(api.ApplicationName, api.ApplicationPassword) error
}

func New(authApplicationPassword api.ApplicationPassword) *applicationDirectory {
	a := new(applicationDirectory)
	a.createAuthApplication(authApplicationPassword)
	return a
}

func (a *applicationDirectory) createAuthApplication(authApplicationPassword api.ApplicationPassword) {
	err := a.Create(api.Application{
		ApplicationName:     AUTH_APPLICATION_NAME,
		ApplicationPassword: authApplicationPassword,
	})
	if err != nil {
		logger.Fatalf("create auth application failed: %v", err)
	}
}

func (a *applicationDirectory) Check(api.ApplicationName, api.ApplicationPassword) error {
	return nil
}

func (a *applicationDirectory) Create(application api.Application) error {
	return nil
}

func (a *applicationDirectory) Delete(applicationName api.ApplicationName) error {
	return nil
}
