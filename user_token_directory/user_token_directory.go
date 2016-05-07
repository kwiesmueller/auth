package user_token_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

const (
	PREFIX = "user"
)

var (
	logger = log.DefaultLogger
)

type userTokenDirectory struct {
	ledis ledis.Set
}

type UserTokenDirectory interface {
	Add(userName api.UserName, authToken api.AuthToken) error
	Exists(userName api.UserName) (bool, error)
	Contains(userName api.UserName, authToken api.AuthToken) (bool, error)
	Remove(userName api.UserName, authToken api.AuthToken) error
}

func New(ledisClient ledis.Set) *userTokenDirectory {
	u := new(userTokenDirectory)
	u.ledis = ledisClient
	return u
}

func createKey(userName api.UserName) string {
	return fmt.Sprintf("%s:%s", PREFIX, userName)
}

func (u *userTokenDirectory) Add(userName api.UserName, authToken api.AuthToken) error {
	logger.Debugf("add user %s", userName)
	key := createKey(userName)
	return u.ledis.SetAdd(key, string(authToken))
}

func (u *userTokenDirectory) Exists(userName api.UserName) (bool, error) {
	logger.Debugf("exists user %s", userName)
	key := createKey(userName)
	return u.ledis.SetExists(key)
}

func (u *userTokenDirectory) Contains(userName api.UserName, authToken api.AuthToken) (bool, error) {
	logger.Debugf("contains user %s token %s", userName, authToken)
	key := createKey(userName)
	return u.ledis.SetContains(key, string(authToken))
}

func (u *userTokenDirectory) Remove(userName api.UserName, authToken api.AuthToken) error {
	logger.Debugf("remove user %s", userName)
	key := createKey(userName)
	return u.ledis.SetRemove(key, string(authToken))
}
