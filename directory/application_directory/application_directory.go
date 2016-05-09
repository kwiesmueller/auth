package application_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const PREFIX = "application"

type directory struct {
	ledis ledis.Kv
}

type ApplicationDirectory interface {
	Create(applicationName api.ApplicationName, applicationPassword api.ApplicationPassword) error
	Delete(applicationName api.ApplicationName) error
	Get(applicationName api.ApplicationName) (*api.ApplicationPassword, error)
}

func New(ledisClient ledis.Kv) *directory {
	a := new(directory)
	a.ledis = ledisClient
	return a
}

func createKey(applicationName api.ApplicationName) string {
	return fmt.Sprintf("%s:%s", PREFIX, applicationName)
}

func (d *directory) Create(applicationName api.ApplicationName, applicationPassword api.ApplicationPassword) error {
	logger.Debugf("create application: %s", applicationName)
	key := createKey(applicationName)
	return d.ledis.Set(key, string(applicationPassword))
}

func (d *directory) Delete(applicationName api.ApplicationName) error {
	logger.Debugf("delete application: %s", applicationName)
	return d.ledis.Del(createKey(applicationName))
}

func (d *directory) Get(applicationName api.ApplicationName) (*api.ApplicationPassword, error) {
	logger.Debugf("get application: %s", applicationName)
	key := createKey(applicationName)
	value, err := d.ledis.Get(key)
	if err != nil {
		return nil, err
	}
	applicationPassword := api.ApplicationPassword(value)
	return &applicationPassword, nil
}
