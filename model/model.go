package model

import (
	"github.com/golang/glog"
	"strings"
)

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

type ApplicationName string

type ApplicationPassword string

type Url string

type Application struct {
	ApplicationName     ApplicationName
	ApplicationPassword ApplicationPassword
}
