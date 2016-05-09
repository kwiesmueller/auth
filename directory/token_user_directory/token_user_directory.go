package token_user_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

type directory struct {
	ledis ledis.Kv
}

type TokenUserDirectory interface {
	Add(authToken api.AuthToken, userName api.UserName) error
	Exists(authToken api.AuthToken) (bool, error)
	Remove(authToken api.AuthToken) error
	FindUserByAuthToken(authToken api.AuthToken) (*api.UserName, error)
}

const PREFIX = "token_user"

var logger = log.DefaultLogger

func New(ledisClient ledis.Kv) *directory {
	u := new(directory)
	u.ledis = ledisClient
	return u
}

func createKey(authToken api.AuthToken) string {
	return fmt.Sprintf("%s:%s", PREFIX, authToken)
}

func (d *directory) Add(authToken api.AuthToken, userName api.UserName) error {
	logger.Debugf("add token %v to user %v", authToken, userName)
	key := createKey(authToken)
	return d.ledis.Set(key, string(userName))
}

func (d *directory) Exists(authToken api.AuthToken) (bool, error) {
	logger.Debugf("exists token %v for user %v", authToken)
	key := createKey(authToken)
	return d.ledis.Exists(key)
}

func (d *directory) Remove(authToken api.AuthToken) error {
	logger.Debugf("remove token %v from user %v", authToken)
	key := createKey(authToken)
	return d.ledis.Del(key)
}

func (d *directory) FindUserByAuthToken(authToken api.AuthToken) (*api.UserName, error) {
	logger.Debugf("find user for token %v", authToken)
	key := createKey(authToken)
	value, err := d.ledis.Get(key)
	if err != nil {
		return nil, err
	}
	userName := api.UserName(value)
	logger.Debugf("found user %v for token %v", userName, authToken)
	return &userName, nil
}
