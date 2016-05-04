package api

type User string

type Group string

type AuthToken string

type ApplicationName string

type ApplicationPassword string

type ApplicationId string

type Request struct {
	ApplicationName     ApplicationName     `json:"applicatonName"`
	ApplicationPassword ApplicationPassword `json:"applicatonPassword"`
	AuthToken           AuthToken           `json:"authToken"`
}

type Response struct {
	User   *User    `json:"user"`
	Groups *[]Group `json:"groups"`
}
