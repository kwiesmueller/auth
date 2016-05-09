package user_token_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

const PREFIX = "user_token"

var logger = log.DefaultLogger

type directory struct {
	ledis ledis.Set
}

type UserTokenDirectory interface {
	Add(userName api.UserName, authToken api.AuthToken) error
	Exists(userName api.UserName) (bool, error)
	Contains(userName api.UserName, authToken api.AuthToken) (bool, error)
	Remove(userName api.UserName, authToken api.AuthToken) error
	Get(userName api.UserName) (*[]api.AuthToken, error)
	Delete(userName api.UserName) error
}

func New(ledisClient ledis.Set) *directory {
	u := new(directory)
	u.ledis = ledisClient
	return u
}

func createKey(userName api.UserName) string {
	return fmt.Sprintf("%s:%s", PREFIX, userName)
}

func (d *directory) Add(userName api.UserName, authToken api.AuthToken) error {
	logger.Debugf("add token %v user %v", authToken, userName)
	key := createKey(userName)
	return d.ledis.SetAdd(key, string(authToken))
}

func (d *directory) Exists(userName api.UserName) (bool, error) {
	logger.Debugf("exists user %v", userName)
	key := createKey(userName)
	return d.ledis.SetExists(key)
}

func (d *directory) Get(userName api.UserName) (*[]api.AuthToken, error) {
	logger.Debugf("get tokens for user %v", userName)
	key := createKey(userName)
	tokens, err := d.ledis.SetGet(key)
	if err != nil {
		return nil, err
	}
	var result []api.AuthToken
	for _, token := range tokens {
		result = append(result, api.AuthToken(token))
	}
	return &result, nil
}

func (d *directory) Contains(userName api.UserName, authToken api.AuthToken) (bool, error) {
	logger.Debugf("contains user %v token %v", userName, authToken)
	key := createKey(userName)
	return d.ledis.SetContains(key, string(authToken))
}

func (d *directory) Remove(userName api.UserName, authToken api.AuthToken) error {
	logger.Debugf("remove token %v from user %v", authToken, userName)
	key := createKey(userName)
	return d.ledis.SetRemove(key, string(authToken))
}

func (d *directory) Delete(userName api.UserName) error {
	logger.Debugf("delete user %v", userName)
	key := createKey(userName)
	return d.ledis.SetClear(key)
}
