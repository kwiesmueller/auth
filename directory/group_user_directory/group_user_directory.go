package group_user_directory

import (
	"github.com/bborbe/auth/api"
	"github.com/bborbe/ledis"
)

type directory struct {
	ledis ledis.Set
}

type GroupUserDirectory interface {
	Add(groupName api.GroupName, userName api.UserName) error
	Remove(groupName api.GroupName, userName api.UserName) error
	Contains(groupName api.GroupName, userName api.UserName) (bool, error)
	Delete(groupName api.GroupName) error
}

func New(ledisClient ledis.Set) *directory {
	d := new(directory)
	d.ledis = ledisClient
	return d
}

func (d *directory) Add(groupName api.GroupName, userName api.UserName) error {
	return nil
}

func (d *directory) Remove(groupName api.GroupName, userName api.UserName) error {
	return nil
}

func (d *directory) Contains(groupName api.GroupName, userName api.UserName) (bool, error) {
	return true, nil
}

func (d *directory) Delete(groupName api.GroupName) error {
	return nil
}
