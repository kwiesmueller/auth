package user_data_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	redis "github.com/bborbe/redis_client"
	"github.com/golang/glog"
)

const PREFIX = "user_data"

type directory struct {
	redis redis.Hash
}

type UserDataDirectory interface {
	Set(userName model.UserName, data map[string]string) error
	SetValue(userName model.UserName, key string, value string) error
	Get(userName model.UserName) (map[string]string, error)
	GetValue(userName model.UserName, key string) (string, error)
	Delete(userName model.UserName) error
	DeleteValue(userName model.UserName, key string) error
}

func New(ledisClient redis.Hash) *directory {
	a := new(directory)
	a.redis = ledisClient
	return a
}

func createKey(userName model.UserName) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, userName)
}

func (d *directory) Set(userName model.UserName, data map[string]string) error {
	glog.V(2).Infof("set %v for user %v %v", data, userName)
	key := createKey(userName)
	for k, v := range data {
		if err := d.redis.HashSet(key, k, v); err != nil {
			return err
		}
	}
	return nil
}

func (d *directory) SetValue(userName model.UserName, field string, value string) error {
	glog.V(2).Infof("set %s=%s for user %v", field, value, userName)
	key := createKey(userName)
	return d.redis.HashSet(key, field, value)
}

func (d *directory) Get(userName model.UserName) (map[string]string, error) {
	glog.V(2).Infof("get data of user %v", userName)
	key := createKey(userName)
	return d.redis.HashGetAll(key)
}

func (d *directory) GetValue(userName model.UserName, field string) (string, error) {
	glog.V(2).Infof("get %s of user %v", field, userName)
	key := createKey(userName)
	return d.redis.HashGet(key, field)
}

func (d *directory) Delete(userName model.UserName) error {
	glog.V(2).Infof("delete data of user %v", userName)
	key := createKey(userName)
	return d.redis.HashClear(key)
}

func (d *directory) DeleteValue(userName model.UserName, field string) error {
	glog.V(2).Infof("delete %s of user %v", field, userName)
	key := createKey(userName)
	return d.redis.HashDel(key, field)
}
