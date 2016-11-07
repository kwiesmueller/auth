package service

import "github.com/bborbe/auth/model"

type ApplicationService interface {
	CreateApplication(applicationName model.ApplicationName) (*model.Application, error)
	CreateApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error)
	VerifyApplicationPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (bool, error)
	GetApplication(applicationName model.ApplicationName) (*model.Application, error)
	DeleteApplication(applicationName model.ApplicationName) error
	ExistsApplication(applicationName model.ApplicationName) (bool, error)
}

type UserService interface {
	DeleteUser(userName model.UserName) error
	DeleteUserWithToken(authToken model.AuthToken) error
	CreateUserWithToken(userName model.UserName, authToken model.AuthToken) error
	AddTokenToUserWithToken(newToken model.AuthToken, userToken model.AuthToken) error
	RemoveTokenFromUserWithToken(newToken model.AuthToken, userToken model.AuthToken) error
	VerifyTokenHasGroups(authToken model.AuthToken, requiredGroupNames []model.GroupName) (*model.UserName, error)
	List() ([]model.UserName, error)
}

type UserDataService interface {
	Set(userName model.UserName, data map[string]string) error
	SetValue(userName model.UserName, key string, value string) error
	Get(userName model.UserName) (map[string]string, error)
	GetValue(userName model.UserName, key string) (string, error)
	Delete(userName model.UserName) error
	DeleteValue(userName model.UserName, key string) error
}

type UserGroupService interface {
	AddUserToGroup(userName model.UserName, groupName model.GroupName) error
	RemoveUserFromGroup(userName model.UserName, groupName model.GroupName) error
}
