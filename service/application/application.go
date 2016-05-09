package application

import "github.com/bborbe/auth/directory/application_directory"

type service struct {
	applicationDirectory application_directory.ApplicationDirectory
}

type Service interface {
}

func New(applicationDirectory application_directory.ApplicationDirectory) *service {
	s := new(service)
	s.applicationDirectory = applicationDirectory
	return s
}
