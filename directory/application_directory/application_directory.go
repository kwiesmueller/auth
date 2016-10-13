package application_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	redis "github.com/bborbe/redis_client"
	"github.com/golang/glog"
)

const PREFIX = "application"

type directory struct {
	redis redis.Kv
}

type ApplicationDirectory interface {
	Create(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) error
	Delete(applicationName model.ApplicationName) error
	Get(applicationName model.ApplicationName) (*model.ApplicationPassword, error)
	Exists(applicationName model.ApplicationName) (bool, error)
}

func New(redisClient redis.Kv) *directory {
	a := new(directory)
	a.redis = redisClient
	return a
}

func createKey(applicationName model.ApplicationName) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, applicationName)
}

func (d *directory) Create(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) error {
	glog.V(2).Infof("create application: %s", applicationName)
	key := createKey(applicationName)
	return d.redis.KvSet(key, string(applicationPassword))
}

func (d *directory) Delete(applicationName model.ApplicationName) error {
	glog.V(2).Infof("delete application: %s", applicationName)
	return d.redis.KvDel(createKey(applicationName))
}

func (d *directory) Exists(applicationName model.ApplicationName) (bool, error) {
	glog.V(2).Infof("exists application: %s", applicationName)
	return d.redis.KvExists(createKey(applicationName))
}

func (d *directory) Get(applicationName model.ApplicationName) (*model.ApplicationPassword, error) {
	glog.V(2).Infof("get application: %s", applicationName)
	key := createKey(applicationName)
	value, err := d.redis.KvGet(key)
	if err != nil {
		return nil, err
	}
	applicationPassword := model.ApplicationPassword(value)
	return &applicationPassword, nil
}
