package user_token_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/ledis"
	"github.com/golang/glog"
)

const PREFIX = "user_token"

type directory struct {
	ledis ledis.Set
}

type UserTokenDirectory interface {
	Add(userName model.UserName, authToken model.AuthToken) error
	Exists(userName model.UserName) (bool, error)
	Contains(userName model.UserName, authToken model.AuthToken) (bool, error)
	Remove(userName model.UserName, authToken model.AuthToken) error
	Get(userName model.UserName) ([]model.AuthToken, error)
	Delete(userName model.UserName) error
	List() ([]model.UserName, error)
}

func New(ledisClient ledis.Set) *directory {
	u := new(directory)
	u.ledis = ledisClient
	return u
}

func createKey(userName model.UserName) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, userName)
}

func (d *directory) Add(userName model.UserName, authToken model.AuthToken) error {
	glog.V(2).Infof("add token %v user %v", authToken, userName)
	key := createKey(userName)
	return d.ledis.SetAdd(key, string(authToken))
}

func (d *directory) Exists(userName model.UserName) (bool, error) {
	glog.V(2).Infof("exists user %v", userName)
	key := createKey(userName)
	return d.ledis.SetExists(key)
}

func (d *directory) Get(userName model.UserName) ([]model.AuthToken, error) {
	glog.V(2).Infof("get tokens for user %v", userName)
	key := createKey(userName)
	tokens, err := d.ledis.SetGet(key)
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
	glog.V(2).Infof("contains user %v token %v", userName, authToken)
	key := createKey(userName)
	return d.ledis.SetContains(key, string(authToken))
}

func (d *directory) Remove(userName model.UserName, authToken model.AuthToken) error {
	glog.V(2).Infof("remove token %v from user %v", authToken, userName)
	key := createKey(userName)
	return d.ledis.SetRemove(key, string(authToken))
}

func (d *directory) Delete(userName model.UserName) error {
	glog.V(2).Infof("delete user %v", userName)
	key := createKey(userName)
	return d.ledis.SetClear(key)
}

func (d *directory) List() ([]model.UserName, error) {
	prefix := fmt.Sprintf("%s:", PREFIX)
	list, err := d.ledis.SetList(prefix)
	if err != nil {
		return nil, err
	}
	var result []model.UserName
	for key := range list {
		name := key[len(prefix):]
		result = append(result, model.UserName(name))
	}
	return result, nil
}
