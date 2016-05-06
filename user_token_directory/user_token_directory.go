package user_token_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/http/bearer"
	"github.com/bborbe/log"
)

type userTokenDirectory struct {
}

type UserTokenDirectory interface {
	FindUserByAuthToken(authToken api.AuthToken) (*api.UserName, error)
	IsUserNotFound(err error) bool
}

var (
	NOT_FOUND = fmt.Errorf("user not found")
	logger    = log.DefaultLogger
)

func New() *userTokenDirectory {
	return new(userTokenDirectory)
}

func (u *userTokenDirectory) FindUserByAuthToken(authToken api.AuthToken) (*api.UserName, error) {
	return findUser(authToken)
}

func findUser(authToken api.AuthToken) (*api.UserName, error) {
	logger.Debugf("find user with auth token: %s", authToken)
	source, id, err := bearer.ParseBearerToken(string(authToken))
	if err != nil {
		return nil, err
	}
	if source == "hipchat" && id == "130647" {
		user := api.UserName("bborbe")
		return &user, nil
	}
	if source == "telegram" && id == "112230768" {
		user := api.UserName("bborbe")
		return &user, nil
	}
	return nil, NOT_FOUND
}

func (u *userTokenDirectory) IsUserNotFound(err error) bool {
	return err == NOT_FOUND
}
