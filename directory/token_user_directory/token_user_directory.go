package token_user_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

const PREFIX = "token_user"

type directory struct {
	ledis ledis.Kv
}

type TokenUserDirectory interface {
	Add(authToken model.AuthToken, userName model.UserName) error
	Exists(authToken model.AuthToken) (bool, error)
	Remove(authToken model.AuthToken) error
	FindUserByAuthToken(authToken model.AuthToken) (*model.UserName, error)
}

var logger = log.DefaultLogger

func New(ledisClient ledis.Kv) *directory {
	u := new(directory)
	u.ledis = ledisClient
	return u
}

func createKey(authToken model.AuthToken) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, authToken)
}

func (d *directory) Add(authToken model.AuthToken, userName model.UserName) error {
	logger.Debugf("add token %v to user %v", authToken, userName)
	key := createKey(authToken)
	return d.ledis.Set(key, string(userName))
}

func (d *directory) Exists(authToken model.AuthToken) (bool, error) {
	logger.Debugf("exists token %v for user %v", authToken)
	key := createKey(authToken)
	return d.ledis.Exists(key)
}

func (d *directory) Remove(authToken model.AuthToken) error {
	logger.Debugf("remove token %v from user %v", authToken)
	key := createKey(authToken)
	return d.ledis.Del(key)
}

func (d *directory) FindUserByAuthToken(authToken model.AuthToken) (*model.UserName, error) {
	logger.Debugf("find user for token %v", authToken)
	key := createKey(authToken)
	value, err := d.ledis.Get(key)
	if err != nil {
		return nil, err
	}
	userName := model.UserName(value)
	logger.Debugf("found user %v for token %v", userName, authToken)
	return &userName, nil
}
