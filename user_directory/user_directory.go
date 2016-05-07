package user_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

type userDirectory struct {
	ledis ledis.Kv
}

type UserDirectory interface {
	Add(userName api.UserName) error
	Exists(userName api.UserName) (bool, error)
	Remove(userName api.UserName) error
}

const (
	PREFIX = "user"
)

var (
	logger = log.DefaultLogger
)

func New(ledisClient ledis.Kv) *userDirectory {
	u := new(userDirectory)
	u.ledis = ledisClient
	return u
}

func createKey(userName api.UserName) string {
	return fmt.Sprintf("%s:%s", PREFIX, userName)
}

func (u *userDirectory) Add(userName api.UserName) error {
	logger.Debugf("add user %s", userName)
	key := createKey(userName)
	return u.ledis.Set(key, string(userName))
}

func (u *userDirectory) Exists(userName api.UserName) (bool, error) {
	logger.Debugf("exists user %s", userName)
	key := createKey(userName)
	return u.ledis.Exists(key)
}

func (u *userDirectory) Remove(userName api.UserName) error {
	logger.Debugf("remove user %s", userName)
	key := createKey(userName)
	return u.ledis.Del(key)
}
