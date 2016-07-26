package user_data_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const PREFIX = "user_data"

type directory struct {
	ledis ledis.Hash
}

type UserDataDirectory interface {
	Set(userName model.UserName, data map[string]string) error
	SetValue(userName model.UserName, key string, value string) error
	Get(userName model.UserName) (map[string]string, error)
	GetValue(userName model.UserName, key string) (string, error)
	Delete(userName model.UserName) error
	DeleteValue(userName model.UserName, key string) error
}

func New(ledisClient ledis.Hash) *directory {
	a := new(directory)
	a.ledis = ledisClient
	return a
}

func createKey(userName model.UserName) string {
	return fmt.Sprintf("%s:%s", PREFIX, userName)
}

func (d *directory) Set(userName model.UserName, data map[string]string) error {
	logger.Debugf("set %v for user %v %v", data, userName)
	key := createKey(userName)
	for k, v := range data {
		if err := d.ledis.HashSet(key, k, v); err != nil {
			return err
		}
	}
	return nil
}

func (d *directory) SetValue(userName model.UserName, field string, value string) error {
	logger.Debugf("set %s=%s for user %v", field, value, userName)
	key := createKey(userName)
	return d.ledis.HashSet(key, field, value)
}

func (d *directory) Get(userName model.UserName) (map[string]string, error) {
	logger.Debugf("get data of user %v", userName)
	key := createKey(userName)
	return d.ledis.HashGetAll(key)
}

func (d *directory) GetValue(userName model.UserName, field string) (string, error) {
	logger.Debugf("get %s of user %v", field, userName)
	key := createKey(userName)
	return d.ledis.HashGet(key, field)
}

func (d *directory) Delete(userName model.UserName) error {
	logger.Debugf("delete data of user %v", userName)
	key := createKey(userName)
	return d.ledis.HashClear(key)
}

func (d *directory) DeleteValue(userName model.UserName, field string) error {
	logger.Debugf("delete %s of user %v", field, userName)
	key := createKey(userName)
	return d.ledis.HashDel(key, field)
}
