package user_group

import (
	"fmt"
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	"github.com/golang/glog"
)

type callRest func(path string, method string, request interface{}, response interface{}) error

type userGroupService struct {
	callRest callRest
}

func New(
	callRest callRest,
) *userGroupService {
	u := new(userGroupService)
	u.callRest = callRest
	return u
}

func (u *userGroupService) AddUserToGroup(userName model.UserName, groupName model.GroupName) error {
	glog.V(4).Infof("add user %s to group %s", userName, groupName)
	request := v1.AddUserToGroupRequest{
		UserName:  model.UserName(userName),
		GroupName: model.GroupName(groupName),
	}
	var response v1.AddUserToGroupResponse
	if err := u.callRest("/api/1.0/user_group", http.MethodPost, &request, &response); err != nil {
		glog.V(2).Infof("add user %v to group %v failed: %v", userName, groupName, err)
		return err
	}
	glog.V(4).Infof("add user user %v to group %v successful", userName, groupName)
	return nil
}

func (u *userGroupService) RemoveUserFromGroup(userName model.UserName, groupName model.GroupName) error {
	glog.V(4).Infof("remove user %s from group %s", userName, groupName)
	request := v1.AddUserToGroupRequest{
		UserName:  model.UserName(userName),
		GroupName: model.GroupName(groupName),
	}
	var response v1.AddUserToGroupResponse
	if err := u.callRest("/api/1.0/user_group", http.MethodDelete, &request, &response); err != nil {
		glog.V(2).Infof("remove user %v from group %v failed: %v", userName, groupName, err)
		return err
	}
	glog.V(4).Infof("remove user user %v from group %v successful", userName, groupName)
	return nil
}

func (u *userGroupService) ListGroupNamesForUsername(username model.UserName) ([]model.GroupName, error) {
	glog.V(4).Infof("list groupnames of user %v", username)
	result := []model.GroupName{}
	if err := u.callRest(fmt.Sprintf("/api/1.0/user_group?username=%v", username), http.MethodGet, nil, &result); err != nil {
		glog.V(2).Infof("list groupnames of user %v failed: %v", username, err)
		return nil, err
	}
	glog.V(4).Infof("got %d groupnames of user %v", len(result), username)
	return result, nil
}
