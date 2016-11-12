package username_token_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	redis "github.com/bborbe/redis_client"
	"github.com/golang/glog"
)

const PREFIX = "user_token"

type directory struct {
	redis redis.Set
}

type UsernameTokenDirectory interface {
	Add(userName model.UserName, authToken model.AuthToken) error
	Exists(userName model.UserName) (bool, error)
	Contains(userName model.UserName, authToken model.AuthToken) (bool, error)
	Remove(userName model.UserName, authToken model.AuthToken) error
	Get(userName model.UserName) ([]model.AuthToken, error)
	Delete(userName model.UserName) error
	List() ([]model.UserName, error)
}

func New(ledisClient redis.Set) *directory {
	u := new(directory)
	u.redis = ledisClient
	return u
}

func createKey(userName model.UserName) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, userName)
}

func (d *directory) Add(userName model.UserName, authToken model.AuthToken) error {
	glog.V(4).Infof("add token %v user %v", authToken, userName)
	key := createKey(userName)
	return d.redis.SetAdd(key, string(authToken))
}

func (d *directory) Exists(userName model.UserName) (bool, error) {
	glog.V(4).Infof("exists user %v", userName)
	key := createKey(userName)
	return d.redis.SetExists(key)
}

func (d *directory) Get(userName model.UserName) ([]model.AuthToken, error) {
	glog.V(4).Infof("get tokens for user %v", userName)
	key := createKey(userName)
	tokens, err := d.redis.SetGet(key)
	if err != nil {
		return nil, err
	}
	var result []model.AuthToken
	for _, token := range tokens {
		result = append(result, model.AuthToken(token))
	}
	return result, nil
}

func (d *directory) Contains(userName model.UserName, authToken model.AuthToken) (bool, error) {
	glog.V(4).Infof("contains user %v token %v", userName, authToken)
	key := createKey(userName)
	return d.redis.SetContains(key, string(authToken))
}

func (d *directory) Remove(userName model.UserName, authToken model.AuthToken) error {
	glog.V(4).Infof("remove token %v from user %v", authToken, userName)
	key := createKey(userName)
	return d.redis.SetRemove(key, string(authToken))
}

func (d *directory) Delete(userName model.UserName) error {
	glog.V(4).Infof("delete user %v", userName)
	key := createKey(userName)
	return d.redis.SetClear(key)
}

func (d *directory) List() ([]model.UserName, error) {
	prefix := fmt.Sprintf("%s:*", PREFIX)
	list, err := d.redis.SetList(prefix)
	if err != nil {
		return nil, err
	}
	var result []model.UserName
	for key := range list {
		name := key[len(prefix)-1:]
		result = append(result, model.UserName(name))
	}
	return result, nil
}
