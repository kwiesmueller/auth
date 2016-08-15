package user_group_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
)

const PREFIX = "user_group"

var logger = log.DefaultLogger

type directory struct {
	ledis ledis.Set
}

type UserGroupDirectory interface {
	Add(userName model.UserName, groupName model.GroupName) error
	Exists(userName model.UserName) (bool, error)
	Get(userName model.UserName) ([]model.GroupName, error)
	Remove(userName model.UserName, groupName model.GroupName) error
	Contains(userName model.UserName, groupName model.GroupName) (bool, error)
	Delete(userName model.UserName) error
}

func New(ledisClient ledis.Set) *directory {
	d := new(directory)
	d.ledis = ledisClient
	return d
}

func createKey(userName model.UserName) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, userName)
}

func (d *directory) Add(userName model.UserName, groupName model.GroupName) error {
	logger.Debugf("add group %v to user %v", groupName, userName)
	key := createKey(userName)
	return d.ledis.SetAdd(key, string(groupName))
}

func (d *directory) Exists(userName model.UserName) (bool, error) {
	logger.Debugf("exists user %v", userName)
	key := createKey(userName)
	return d.ledis.SetExists(key)
}

func (d *directory) Get(userName model.UserName) ([]model.GroupName, error) {
	logger.Debugf("get groups of user %v", userName)
	key := createKey(userName)
	groups, err := d.ledis.SetGet(key)
	if err != nil {
		return nil, err
	}
	var result []model.GroupName
	for _, group := range groups {
		result = append(result, model.GroupName(group))
	}
	return result, nil
}

func (d *directory) Remove(userName model.UserName, groupName model.GroupName) error {
	logger.Debugf("remove group %v from user %v", groupName, groupName)
	key := createKey(userName)
	return d.ledis.SetRemove(key, string(groupName))
}

func (d *directory) Contains(userName model.UserName, groupName model.GroupName) (bool, error) {
	logger.Debugf("contains user %v group %v", userName, groupName)
	key := createKey(userName)
	return d.ledis.SetContains(key, string(groupName))
}

func (d *directory) Delete(userName model.UserName) error {
	logger.Debugf("delete user %v", userName)
	key := createKey(userName)
	return d.ledis.SetClear(key)
}
