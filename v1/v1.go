package v1

import "github.com/bborbe/auth/model"

const VERSION = "1.0"

type User struct {
	UserName model.UserName `json:"username"`
}

type LoginRequest struct {
	AuthToken      model.AuthToken   `json:"authToken"`
	RequiredGroups []model.GroupName `json:"groups"`
}

type LoginResponse struct {
	UserName *model.UserName `json:"user"`
}

type CreateApplicationRequest struct {
	ApplicationName model.ApplicationName `json:"applicatonName"`
}

type CreateApplicationResponse struct {
	ApplicationName     model.ApplicationName     `json:"applicatonName"`
	ApplicationPassword model.ApplicationPassword `json:"applicatonPassword"`
}

type DeleteApplicationRequest struct {
}

type DeleteApplicationResponse struct {
}

type GetApplicationRequest struct {
}

type GetApplicationResponse struct {
	ApplicationName     model.ApplicationName     `json:"applicatonName"`
	ApplicationPassword model.ApplicationPassword `json:"applicatonPassword"`
}

type RegisterRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
	UserName  model.UserName  `json:"user"`
}

type RegisterResponse struct {
}

type UnregisterRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
}

type UnregisterResponse struct {
}

type AddTokenRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
	Token     model.AuthToken `json:"token"`
}

type AddTokenResponse struct {
}

type RemoveTokenRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
	Token     model.AuthToken `json:"token"`
}

type RemoveTokenResponse struct {
}

type AddUserToGroupRequest struct {
	UserName  model.UserName  `json:"user"`
	GroupName model.GroupName `json:"group"`
}

type AddUserToGroupResponse struct {
}

type RemoveUserFromGroupRequest struct {
	UserName  model.UserName  `json:"user"`
	GroupName model.GroupName `json:"group"`
}

type RemoveUserFromGroupResponse struct {
}

type SetUserDataRequest map[string]string

type SetUserDataResponse struct {
}

type SetUserDataValueRequest string

type SetUserDataValueResponse struct {
}

type GetUserDataRequest struct {
}

type GetUserDataResponse map[string]string

type GetUserDataValueRequest struct {
}

type GetUserDataValueResponse string

type DeleteUserDataRequest struct {
}

type DeleteUserDataResponse struct {
}

type DeleteUserDataValueRequest struct {
}

type DeleteUserDataValueResponse struct {
}

type UserListRequest struct {
}

type UserListResponse []User
