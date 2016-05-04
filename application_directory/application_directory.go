package application_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/ledis"
	"github.com/timelinelabs/romulus/Godeps/_workspace/src/github.com/albertrdixon/gearbox/logger"
)

const (
	AUTH_APPLICATION_NAME = api.ApplicationName("auth")
	PREFIX                = "application"
	FIELD_PASSWORD        = "password"
)

type applicationDirectory struct {
	ledis ledis.Client
}

type ApplicationDirectory interface {
	Check(api.ApplicationName, api.ApplicationPassword) error
	Create(application api.Application) error
	Delete(applicationName api.ApplicationName) error
}

func New(ledisClient ledis.Client) *applicationDirectory {
	a := new(applicationDirectory)
	a.ledis = ledisClient
	return a
}

func (a *applicationDirectory) Check(applicationName api.ApplicationName, applicationPassword api.ApplicationPassword) error {
	logger.Debugf("check application: %s", applicationName)
	value, err := a.ledis.HashGet(createKey(applicationName), FIELD_PASSWORD)
	if err != nil {
		return err
	}
	if api.ApplicationPassword(value) != applicationPassword {
		return fmt.Errorf("password invalid")
	}
	return nil
}

func (a *applicationDirectory) Create(application api.Application) error {
	logger.Debugf("create application: %s", application.ApplicationName)
	key := createKey(application.ApplicationName)
	return a.ledis.HashSet(key, FIELD_PASSWORD, string(application.ApplicationPassword))
}

func (a *applicationDirectory) Delete(applicationName api.ApplicationName) error {
	logger.Debugf("delete application: %s", applicationName)
	return a.ledis.HashClear(createKey(applicationName))
}

func createKey(applicationName api.ApplicationName) string {
	return fmt.Sprintf("%s:%s", PREFIX, applicationName)
}
