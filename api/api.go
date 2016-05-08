package api

type UserName string

type GroupName string

type AuthToken string

type ApplicationName string

type ApplicationPassword string

type Application struct {
	ApplicationName     ApplicationName
	ApplicationPassword ApplicationPassword
}

type LoginRequest struct {
	AuthToken      AuthToken   `json:"authToken"`
	RequiredGroups []GroupName `json:"groups"`
}

type LoginResponse struct {
	UserName *UserName `json:"user"`
}

type CreateApplicationRequest struct {
	ApplicationName ApplicationName `json:"applicatonName"`
}

type CreateApplicationResponse struct {
	ApplicationName     ApplicationName     `json:"applicatonName"`
	ApplicationPassword ApplicationPassword `json:"applicatonPassword"`
}

type DeleteApplicationRequest struct {
}

type DeleteApplicationResponse struct {
}

type GetApplicationRequest struct {
}

type GetApplicationResponse struct {
	ApplicationName     ApplicationName     `json:"applicatonName"`
	ApplicationPassword ApplicationPassword `json:"applicatonPassword"`
}

type RegisterRequest struct {
	AuthToken AuthToken `json:"authToken"`
	UserName  UserName  `json:"user"`
}

type RegisterResponse struct {
}

type UnRegisterRequest struct {
	AuthToken AuthToken `json:"authToken"`
}

type UnRegisterResponse struct {
}

type AddTokenRequest struct {
	AuthToken AuthToken `json:"authToken"`
	Token     AuthToken `json:"token"`
}

type AddTokenResponse struct {
}

type RemoveTokenRequest struct {
	AuthToken AuthToken `json:"authToken"`
	Token     AuthToken `json:"token"`
}

type RemoveTokenResponse struct {
}
