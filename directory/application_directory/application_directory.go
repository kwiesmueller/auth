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
	PREFIX         = "application"
	FIELD_PASSWORD = "password"
)

type directory struct {
	ledis ledis.Hash
}

type ApplicationDirectory interface {
	Check(api.ApplicationName, api.ApplicationPassword) error
	Create(application api.Application) error
	Delete(applicationName api.ApplicationName) error
	Get(applicationName api.ApplicationName) (*api.Application, error)
	IsApplicationNotFound(err error) bool
}

func New(ledisClient ledis.Hash) *directory {
	a := new(directory)
	a.ledis = ledisClient
	return a
}

func (d *directory) IsApplicationNotFound(err error) bool {
	return err == NOT_FOUND
}

func (d *directory) Check(applicationName api.ApplicationName, applicationPassword api.ApplicationPassword) error {
	logger.Debugf("check application: %s", applicationName)
	value, err := d.ledis.HashGet(createKey(applicationName), FIELD_PASSWORD)
	if err != nil {
		return err
	}
	if api.ApplicationPassword(value) != applicationPassword {
		return fmt.Errorf("password invalid")
	}
	return nil
}

func (d *directory) Create(application api.Application) error {
	logger.Debugf("create application: %s", application.ApplicationName)
	key := createKey(application.ApplicationName)
	return d.ledis.HashSet(key, FIELD_PASSWORD, string(application.ApplicationPassword))
}

func (d *directory) Delete(applicationName api.ApplicationName) error {
	logger.Debugf("delete application: %s", applicationName)
	return d.ledis.HashClear(createKey(applicationName))
}

func createKey(applicationName api.ApplicationName) string {
	return fmt.Sprintf("%s:%s", PREFIX, applicationName)
}

func (d *directory) Get(applicationName api.ApplicationName) (*api.Application, error) {
	logger.Debugf("get application: %s", applicationName)
	key := createKey(applicationName)
	exists, err := d.ledis.HashExists(key, FIELD_PASSWORD)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, NOT_FOUND
	}
	value, err := d.ledis.HashGet(key, FIELD_PASSWORD)
	if err != nil {
		return nil, err
	}
	return &api.Application{
		ApplicationName:     applicationName,
		ApplicationPassword: api.ApplicationPassword(value),
	}, nil
}
