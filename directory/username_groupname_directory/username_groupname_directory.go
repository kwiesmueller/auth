package username_groupname_directory

import (
	"fmt"

	"github.com/bborbe/auth/model"
	redis "github.com/bborbe/redis_client"
	"github.com/golang/glog"
)

const PREFIX = "user_group"

type directory struct {
	redis redis.Set
}

type UsernameGroupnameDirectory interface {
	Add(userName model.UserName, groupName model.GroupName) error
	Exists(userName model.UserName) (bool, error)
	Get(userName model.UserName) ([]model.GroupName, error)
	Remove(userName model.UserName, groupName model.GroupName) error
	Contains(userName model.UserName, groupName model.GroupName) (bool, error)
	Delete(userName model.UserName) error
}

func New(ledisClient redis.Set) *directory {
	d := new(directory)
	d.redis = ledisClient
	return d
}

func createKey(userName model.UserName) string {
	return fmt.Sprintf("%s%s%s", PREFIX, model.Seperator, userName)
}

func (d *directory) Add(userName model.UserName, groupName model.GroupName) error {
	glog.V(2).Infof("add group %v to user %v", groupName, userName)
	key := createKey(userName)
	return d.redis.SetAdd(key, string(groupName))
}

func (d *directory) Exists(userName model.UserName) (bool, error) {
	glog.V(2).Infof("exists user %v", userName)
	key := createKey(userName)
	return d.redis.SetExists(key)
}

func (d *directory) Get(userName model.UserName) ([]model.GroupName, error) {
	glog.V(2).Infof("get groups of user %v", userName)
	key := createKey(userName)
	groups, err := d.redis.SetGet(key)
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
	glog.V(2).Infof("remove group %v from user %v", groupName, groupName)
	key := createKey(userName)
	return d.redis.SetRemove(key, string(groupName))
}

func (d *directory) Contains(userName model.UserName, groupName model.GroupName) (bool, error) {
	glog.V(2).Infof("contains user %v group %v", userName, groupName)
	key := createKey(userName)
	return d.redis.SetContains(key, string(groupName))
}

func (d *directory) Delete(userName model.UserName) error {
	glog.V(2).Infof("delete user %v", userName)
	key := createKey(userName)
	return d.redis.SetClear(key)
}
