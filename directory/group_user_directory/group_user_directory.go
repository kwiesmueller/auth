package group_user_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	redis "github.com/bborbe/redis_client"
	"github.com/golang/glog"
)

const PREFIX = "group_user"

type directory struct {
	redis redis.Set
}

type GroupUserDirectory interface {
	Add(groupName model.GroupName, userName model.UserName) error
	Exists(groupName model.GroupName) (bool, error)
	Get(groupName model.GroupName) ([]model.UserName, error)
	Remove(groupName model.GroupName, userName model.UserName) error
	Contains(groupName model.GroupName, userName model.UserName) (bool, error)
	Delete(groupName model.GroupName) error
}

func New(ledisClient redis.Set) *directory {
	d := new(directory)
	d.redis = ledisClient
	return d
}

func createKey(groupName model.GroupName) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, groupName)
}

func (d *directory) Add(groupName model.GroupName, userName model.UserName) error {
	glog.V(4).Infof("add user %v to group %v", userName, groupName)
	key := createKey(groupName)
	return d.redis.SetAdd(key, string(userName))
}

func (d *directory) Exists(groupName model.GroupName) (bool, error) {
	glog.V(4).Infof("exists group %v", groupName)
	key := createKey(groupName)
	return d.redis.SetExists(key)
}

func (d *directory) Get(groupName model.GroupName) ([]model.UserName, error) {
	glog.V(4).Infof("get users of group %v", groupName)
	key := createKey(groupName)
	users, err := d.redis.SetGet(key)
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
	glog.V(4).Infof("remove user %v from group %v", userName, userName)
	key := createKey(groupName)
	return d.redis.SetRemove(key, string(userName))
}

func (d *directory) Contains(groupName model.GroupName, userName model.UserName) (bool, error) {
	glog.V(4).Infof("contains group %v user %v", groupName, userName)
	key := createKey(groupName)
	return d.redis.SetContains(key, string(userName))
}

func (d *directory) Delete(groupName model.GroupName) error {
	glog.V(4).Infof("delete group %v", groupName)
	key := createKey(groupName)
	return d.redis.SetClear(key)
}
