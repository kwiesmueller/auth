package group_user_directory

import (
	"fmt"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

const PREFIX = "group_user"

var logger = log.DefaultLogger

type directory struct {
	ledis ledis.Set
}

type GroupUserDirectory interface {
	Add(groupName api.GroupName, userName api.UserName) error
	Exists(groupName api.GroupName) (bool, error)
	Get(groupName api.GroupName) (*[]api.UserName, error)
	Remove(groupName api.GroupName, userName api.UserName) error
	Contains(groupName api.GroupName, userName api.UserName) (bool, error)
	Delete(groupName api.GroupName) error
}

func New(ledisClient ledis.Set) *directory {
	d := new(directory)
	d.ledis = ledisClient
	return d
}

func createKey(groupName api.GroupName) string {
	return fmt.Sprintf("%s:%s", PREFIX, groupName)
}

func (d *directory) Add(groupName api.GroupName, userName api.UserName) error {
	logger.Debugf("add user %v to group %v", userName, groupName)
	key := createKey(groupName)
	return d.ledis.SetAdd(key, string(userName))
}

func (d *directory) Exists(groupName api.GroupName) (bool, error) {
	logger.Debugf("exists group %v", groupName)
	key := createKey(groupName)
	return d.ledis.SetExists(key)
}

func (d *directory) Get(groupName api.GroupName) (*[]api.UserName, error) {
	logger.Debugf("get users of group %v", groupName)
	key := createKey(groupName)
	users, err := d.ledis.SetGet(key)
	if err != nil {
		return nil, err
	}
	var result []api.UserName
	for _, user := range users {
		result = append(result, api.UserName(user))
	}
	return &result, nil
}

func (d *directory) Remove(groupName api.GroupName, userName api.UserName) error {
	logger.Debugf("remove user %v from group %v", userName, userName)
	key := createKey(groupName)
	return d.ledis.SetRemove(key, string(userName))
}

func (d *directory) Contains(groupName api.GroupName, userName api.UserName) (bool, error) {
	logger.Debugf("contains group %v user %v", groupName, userName)
	key := createKey(groupName)
	return d.ledis.SetContains(key, string(userName))
}

func (d *directory) Delete(groupName api.GroupName) error {
	logger.Debugf("delete group %v", groupName)
	key := createKey(groupName)
	return d.ledis.SetClear(key)
}
