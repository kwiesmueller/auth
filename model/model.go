package model

const (
	AUTH_APPLICATION_NAME = ApplicationName("auth")
	AUTH_ADMIN_GROUP      = GroupName("auth")
)

type UserName string

type GroupName string

type AuthToken string

type ApplicationName string

type ApplicationPassword string

type Address string

type Application struct {
	ApplicationName     ApplicationName
	ApplicationPassword ApplicationPassword
}
