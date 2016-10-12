package token_user_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	redis "github.com/bborbe/redis_client"
	"github.com/golang/glog"
)

const PREFIX = "token_user"

type directory struct {
	redis redis.Kv
}

type TokenUserDirectory interface {
	Add(authToken model.AuthToken, userName model.UserName) error
	Exists(authToken model.AuthToken) (bool, error)
	Remove(authToken model.AuthToken) error
	FindUserByAuthToken(authToken model.AuthToken) (*model.UserName, error)
}

func New(ledisClient redis.Kv) *directory {
	u := new(directory)
	u.redis = ledisClient
	return u
}

func createKey(authToken model.AuthToken) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, authToken)
}

func (d *directory) Add(authToken model.AuthToken, userName model.UserName) error {
	glog.V(2).Infof("add token %v to user %v", authToken, userName)
	key := createKey(authToken)
	return d.redis.Set(key, string(userName))
}

func (d *directory) Exists(authToken model.AuthToken) (bool, error) {
	glog.V(2).Infof("exists token %v", authToken)
	key := createKey(authToken)
	return d.redis.Exists(key)
}

func (d *directory) Remove(authToken model.AuthToken) error {
	glog.V(2).Infof("remove token %v", authToken)
	key := createKey(authToken)
	return d.redis.Del(key)
}

func (d *directory) FindUserByAuthToken(authToken model.AuthToken) (*model.UserName, error) {
	glog.V(2).Infof("find user for token %v", authToken)
	key := createKey(authToken)
	value, err := d.redis.Get(key)
	if err != nil {
		return nil, err
	}
	userName := model.UserName(value)
	glog.V(2).Infof("found user %v for token %v", userName, authToken)
	return &userName, nil
}
