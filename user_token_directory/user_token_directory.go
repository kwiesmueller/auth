package user_token_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/ledis"
	"github.com/bborbe/log"
)

type userTokenDirectory struct {
	ledis ledis.Client
}

type UserTokenDirectory interface {
	FindUserByAuthToken(authToken api.AuthToken) (*api.UserName, error)
	IsUserNotFound(err error) bool
	Add(userName api.UserName, authToken api.AuthToken) error
	Remove(userName api.UserName, authToken api.AuthToken) error
}

const (
	PREFIX = "application"
)

var (
	NOT_FOUND = fmt.Errorf("user not found")
	logger    = log.DefaultLogger
)

func New(ledisClient ledis.Client) *userTokenDirectory {
	u := new(userTokenDirectory)
	u.ledis = ledisClient
	return u
}

func createKey(authToken api.AuthToken) string {
	return fmt.Sprintf("%s:%s", PREFIX, authToken)
}

func (u *userTokenDirectory) Add(userName api.UserName, authToken api.AuthToken) error {
	logger.Debugf("add token %s to user %s", authToken, userName)
	key := createKey(authToken)
	return u.ledis.Set(key, string(userName))
}

func (u *userTokenDirectory) Remove(userName api.UserName, authToken api.AuthToken) error {
	logger.Debugf("remove token %s from user %s", authToken, userName)
	key := createKey(authToken)
	return u.ledis.Del(key)
}

func (u *userTokenDirectory) FindUserByAuthToken(authToken api.AuthToken) (*api.UserName, error) {
	logger.Debugf("find user for token %s", authToken)
	key := createKey(authToken)
	value, err := u.ledis.Get(key)
	if err != nil {
		return nil, err
	}
	userName := api.UserName(value)
	logger.Debugf("found user %s for token %s", userName, authToken)
	return &userName, nil
}

func (u *userTokenDirectory) IsUserNotFound(err error) bool {
	return err == NOT_FOUND
}
