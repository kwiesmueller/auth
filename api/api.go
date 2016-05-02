package api

type User string

type Group string

type AuthToken string

type Request struct {
	ApplicationName     string    `json:"applicatonName"`
	ApplicationPassword string    `json:"applicatonPassword"`
	AuthToken           AuthToken `json:"authToken"`
}

type Response struct {
	User   *User    `json:"user"`
	Groups *[]Group `json:"groups"`
}
