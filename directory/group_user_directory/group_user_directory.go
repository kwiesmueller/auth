package group_user_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

const PREFIX = "group_user"

var logger = log.DefaultLogger

type directory struct {
	ledis ledis.Set
}

type GroupUserDirectory interface {
	Add(groupName model.GroupName, userName model.UserName) error
	Exists(groupName model.GroupName) (bool, error)
	Get(groupName model.GroupName) ([]model.UserName, error)
	Remove(groupName model.GroupName, userName model.UserName) error
	Contains(groupName model.GroupName, userName model.UserName) (bool, error)
	Delete(groupName model.GroupName) error
}

func New(ledisClient ledis.Set) *directory {
	d := new(directory)
	d.ledis = ledisClient
	return d
}

func createKey(groupName model.GroupName) string {
	return fmt.Sprintf("%s:%s", PREFIX, groupName)
}

func (d *directory) Add(groupName model.GroupName, userName model.UserName) error {
	logger.Debugf("add user %v to group %v", userName, groupName)
	key := createKey(groupName)
	return d.ledis.SetAdd(key, string(userName))
}

func (d *directory) Exists(groupName model.GroupName) (bool, error) {
	logger.Debugf("exists group %v", groupName)
	key := createKey(groupName)
	return d.ledis.SetExists(key)
}

func (d *directory) Get(groupName model.GroupName) ([]model.UserName, error) {
	logger.Debugf("get users of group %v", groupName)
	key := createKey(groupName)
	users, err := d.ledis.SetGet(key)
	if err != nil {
		return nil, err
	}
	var result []model.UserName
	for _, user := range users {
		result = append(result, model.UserName(user))
	}
	return result, nil
}

func (d *directory) Remove(groupName model.GroupName, userName model.UserName) error {
	logger.Debugf("remove user %v from group %v", userName, userName)
	key := createKey(groupName)
	return d.ledis.SetRemove(key, string(userName))
}

func (d *directory) Contains(groupName model.GroupName, userName model.UserName) (bool, error) {
	logger.Debugf("contains group %v user %v", groupName, userName)
	key := createKey(groupName)
	return d.ledis.SetContains(key, string(userName))
}

func (d *directory) Delete(groupName model.GroupName) error {
	logger.Debugf("delete group %v", groupName)
	key := createKey(groupName)
	return d.ledis.SetClear(key)
}
