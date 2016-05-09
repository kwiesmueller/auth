package application

import (
	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/log"
)

const PASSWORD_LENGTH = 16

var logger = log.DefaultLogger

type GeneratePassword func(length int) string

type service struct {
	generatePassword     GeneratePassword
	applicationDirectory application_directory.ApplicationDirectory
}

type Service interface {
	CreateApplication(applicationName api.ApplicationName) (*api.Application, error)
}

func New(generatePassword GeneratePassword, applicationDirectory application_directory.ApplicationDirectory) *service {
	s := new(service)
	s.generatePassword = generatePassword
	s.applicationDirectory = applicationDirectory
	return s
}

func (s *service) CreateApplication(applicationName api.ApplicationName) (*api.Application, error) {
	logger.Debugf("create application %v", applicationName)
	application := api.Application{
		ApplicationName:     applicationName,
		ApplicationPassword: api.ApplicationPassword(s.generatePassword(PASSWORD_LENGTH)),
	}
	if err := s.applicationDirectory.Create(application); err != nil {
		logger.Debugf("create application %v failed: %v", applicationName, err)
		return nil, err
	}
	logger.Debugf("created application %v successful", applicationName)
	return &application, nil
}
