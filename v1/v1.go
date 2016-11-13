package v1

import "github.com/bborbe/auth/model"

const VERSION = "1.0"

type LoginRequest struct {
	AuthToken      model.AuthToken   `json:"authToken"`
	RequiredGroups []model.GroupName `json:"groups"`
}

type LoginResponse struct {
	UserName *model.UserName `json:"user"`
}

type RegisterRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
	UserName  model.UserName  `json:"user"`
}

type UnregisterRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
}

type AddTokenRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
	Token     model.AuthToken `json:"token"`
}

type RemoveTokenRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
	Token     model.AuthToken `json:"token"`
}

type UsernameTokenRequest struct {
	AuthToken model.AuthToken `json:"authToken"`
	Userame   model.UserName  `json:"username"`
}

type AddUserToGroupRequest struct {
	UserName  model.UserName  `json:"user"`
	GroupName model.GroupName `json:"group"`
}

type RemoveUserFromGroupRequest struct {
	UserName  model.UserName  `json:"user"`
	GroupName model.GroupName `json:"group"`
}
