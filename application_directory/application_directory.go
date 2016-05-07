package application_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

var (
	logger    = log.DefaultLogger
	NOT_FOUND = fmt.Errorf("application not found")
)

const (
	AUTH_APPLICATION_NAME = api.ApplicationName("auth")
	PREFIX                = "application"
	FIELD_PASSWORD        = "password"
)

type applicationDirectory struct {
	ledis ledis.Hash
}

type ApplicationDirectory interface {
	Check(api.ApplicationName, api.ApplicationPassword) error
	Create(application api.Application) error
	Delete(applicationName api.ApplicationName) error
	Get(applicationName api.ApplicationName) (*api.Application, error)
	IsApplicationNotFound(err error) bool
}

func New(ledisClient ledis.Hash) *applicationDirectory {
	a := new(applicationDirectory)
	a.ledis = ledisClient
	return a
}

func (a *applicationDirectory) IsApplicationNotFound(err error) bool {
	return err == NOT_FOUND
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

func (a *applicationDirectory) Get(applicationName api.ApplicationName) (*api.Application, error) {
	logger.Debugf("get application: %s", applicationName)
	key := createKey(applicationName)
	exists, err := a.ledis.HashExists(key, FIELD_PASSWORD)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, NOT_FOUND
	}
	value, err := a.ledis.HashGet(key, FIELD_PASSWORD)
	if err != nil {
		return nil, err
	}
	return &api.Application{
		ApplicationName:     applicationName,
		ApplicationPassword: api.ApplicationPassword(value),
	}, nil
}
