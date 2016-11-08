package model

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// Port to listen on
type Port int

// Address representation of the port
func (p Port) Address() string {
	return fmt.Sprintf(":%d", p)
}

// Int value of the port
func (p Port) Int() int {
	return int(p)
}

// Prefix of the application
type Prefix string

// String represenation of the prefix
func (p Prefix) String() string {
	return string(p)
}

// LedisdbAddress is used to connect to ledis (localhost:5555)
type LedisdbAddress string

func (l LedisdbAddress) String() string {
	return string(l)
}

// LedisdbPassword used to access ledis
type LedisdbPassword string

func (l LedisdbPassword) String() string {
	return string(l)
}

const (
	Seperator             = ":"
	AUTH_APPLICATION_NAME = ApplicationName("auth")
	AUTH_ADMIN_GROUP      = GroupName("auth")
)

type UserName string

type GroupName string

func CreateGroupsFromString(groupNames string) []GroupName {
	parts := strings.Split(groupNames, ",")
	groups := make([]GroupName, 0)
	for _, groupName := range parts {
		if len(groupName) > 0 {
			groups = append(groups, GroupName(groupName))
		}
	}
	glog.V(1).Infof("required groups: %v", groups)
	return groups
}

type AuthToken string

func (a AuthToken) String() string {
	return string(a)
}

func ParseTokens(tokens string) []AuthToken {
	var result []AuthToken
	for _, token := range strings.Split(tokens, ",") {
		if len(token) > 0 {
			result = append(result, AuthToken(token))
		}
	}
	return result
}

type ApplicationName string

func (a ApplicationName) String() string {
	return string(a)
}

type ApplicationPassword string

func (a ApplicationPassword) String() string {
	return string(a)
}

type Url string

func (u Url) String() string {
	return string(u)
}

type Application struct {
	ApplicationName     ApplicationName
	ApplicationPassword ApplicationPassword
}
