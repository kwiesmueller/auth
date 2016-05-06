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
	AuthToken AuthToken `json:"authToken"`
}

type LoginResponse struct {
	User   *UserName    `json:"user"`
	Groups *[]GroupName `json:"groups"`
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
