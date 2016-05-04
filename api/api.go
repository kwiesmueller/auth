package api

type User string

type Group string

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
	User   *User    `json:"user"`
	Groups *[]Group `json:"groups"`
}

type CreateApplicationRequest struct {
	ApplicationName     ApplicationName     `json:"applicatonName"`
	ApplicationPassword ApplicationPassword `json:"applicatonPassword"`
	AuthToken           AuthToken           `json:"authToken"`
}
