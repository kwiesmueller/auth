package user_group_directory

import (
	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
)

type directory struct {
	ledis ledis.Set
}

type UserGroupDirectory interface {
	Add(userName api.UserName, groupName api.GroupName) error
	Remove(userName api.UserName, groupName api.GroupName) error
	Contains(userName api.UserName, groupName api.GroupName) (bool, error)
	Delete(userName api.UserName) error
}

func New(ledisClient ledis.Set) *directory {
	d := new(directory)
	d.ledis = ledisClient
	return d
}

func (d *directory) Add(userName api.UserName, groupName api.GroupName) error {
	return nil
}

func (d *directory) Remove(userName api.UserName, groupName api.GroupName) error {
	return nil
}

func (d *directory) Contains(userName api.UserName, groupName api.GroupName) (bool, error) {
	return true, nil
}

func (d *directory) Delete(userName api.UserName) error {
	return nil
}
