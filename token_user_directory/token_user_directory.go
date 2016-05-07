package token_user_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

type tokenUserDirectory struct {
	ledis ledis.Kv
}

type TokenUserDirectory interface {
	Add(authToken api.AuthToken, userName api.UserName) error
	Exists(authToken api.AuthToken) (bool, error)
	Remove(authToken api.AuthToken, userName api.UserName) error
	FindUserByAuthToken(authToken api.AuthToken) (*api.UserName, error)
	IsUserNotFound(err error) bool
}

const (
	PREFIX = "token_user"
)

var (
	NOT_FOUND = fmt.Errorf("user not found")
	logger    = log.DefaultLogger
)

func New(ledisClient ledis.Kv) *tokenUserDirectory {
	u := new(tokenUserDirectory)
	u.ledis = ledisClient
	return u
}

func createKey(authToken api.AuthToken) string {
	return fmt.Sprintf("%s:%s", PREFIX, authToken)
}

func (u *tokenUserDirectory) Add(authToken api.AuthToken, userName api.UserName) error {
	logger.Debugf("add token %s to user %s", authToken, userName)
	key := createKey(authToken)
	return u.ledis.Set(key, string(userName))
}

func (u *tokenUserDirectory) Exists(authToken api.AuthToken) (bool, error) {
	logger.Debugf("exists token %s to user %s", authToken)
	key := createKey(authToken)
	return u.ledis.Exists(key)
}

func (u *tokenUserDirectory) Remove(authToken api.AuthToken, userName api.UserName) error {
	logger.Debugf("remove token %s from user %s", authToken, userName)
	key := createKey(authToken)
	return u.ledis.Del(key)
}

func (u *tokenUserDirectory) FindUserByAuthToken(authToken api.AuthToken) (*api.UserName, error) {
	logger.Debugf("find user for token %s", authToken)
	key := createKey(authToken)
	value, err := u.ledis.Get(key)
	if err != nil {
		return nil, NOT_FOUND
	}
	userName := api.UserName(value)
	logger.Debugf("found user %s for token %s", userName, authToken)
	return &userName, nil
}

func (u *tokenUserDirectory) IsUserNotFound(err error) bool {
	return err == NOT_FOUND
}
