package user_data_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const PREFIX = "user_data"

type directory struct {
	ledis ledis.Kv
}

type UserDataDirectory interface {
	SetValue(userName api.UserName, key string, value string) error
	GetValue(userName api.UserName, key string) (string, error)
	Get(userName api.UserName) (map[string]string, error)
	Set(userName api.UserName, data map[string]string) error
	DeleteValue(userName api.UserName, key string) error
	Delete(userName api.UserName) error
}

func New(ledisClient ledis.Kv) *directory {
	a := new(directory)
	a.ledis = ledisClient
	return a
}

func createKey(applicationName api.ApplicationName) string {
	return fmt.Sprintf("%s:%s", PREFIX, applicationName)
}

func (d *directory) Set(userName api.UserName, data map[string]string) error {
	logger.Debugf("set %v for user %v %v", data, userName)
	return nil
}

func (d *directory) SetValue(userName api.UserName, key string, value string) error {
	logger.Debugf("set %s=%s for user %v", key, value, userName)
	return nil
}

func (d *directory) Get(userName api.UserName) (map[string]string, error) {
	logger.Debugf("get data of user %v", userName)
	return nil, nil
}

func (d *directory) GetValue(userName api.UserName, key string) (string, error) {
	logger.Debugf("get %s of user %v", key, userName)
	return "", nil
}

func (d *directory) Delete(userName api.UserName) error {
	logger.Debugf("delete data of user %v", userName)
	return nil
}

func (d *directory) DeleteValue(userName api.UserName, key string) error {
	logger.Debugf("delete %s of user %v", key, userName)
	return nil
}
