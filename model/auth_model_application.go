package model

type Application struct {
	ApplicationName     ApplicationName     `json:"applicatonName"`
	ApplicationPassword ApplicationPassword `json:"applicatonPassword"`
}

type ApplicationName string

func (a ApplicationName) String() string {
	return string(a)
}

type ApplicationPassword string

func (a ApplicationPassword) String() string {
	return string(a)
}
