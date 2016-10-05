package group_checker

import (
	auth_model "github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

type verifyGroups func(authToken auth_model.AuthToken, requiredGroups []auth_model.GroupName) (*auth_model.UserName, error)

type GroupChecker interface {
	HasRequiredGroups(authToken auth_model.AuthToken) bool
}

type groupChecker struct {
	verifyGroups   verifyGroups
	requiredGroups []auth_model.GroupName
}

func New(
	verifyGroups verifyGroups,
	requiredGroups ...auth_model.GroupName,
) *groupChecker {
	g := new(groupChecker)
	g.requiredGroups = requiredGroups
	g.verifyGroups = verifyGroups
	return g
}

func (g *groupChecker) HasRequiredGroups(authToken auth_model.AuthToken) bool {
	glog.V(2).Infof("has required groups")
	username, err := g.verifyGroups(authToken, g.requiredGroups)
	return err == nil && username != nil && len(*username) > 0
}
