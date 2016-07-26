package application

import (
	"fmt"

	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/auth/model"
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
	CreateApplication(applicationName model.ApplicationName) (*model.Application, error)
	CreateApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error)
	VerifyApplicationPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (bool, error)
	GetApplication(applicationName model.ApplicationName) (*model.Application, error)
	DeleteApplication(applicationName model.ApplicationName) error
	ExistsApplication(applicationName model.ApplicationName) (bool, error)
}

func New(generatePassword GeneratePassword, applicationDirectory application_directory.ApplicationDirectory) *service {
	s := new(service)
	s.generatePassword = generatePassword
	s.applicationDirectory = applicationDirectory
	return s
}

func (s *service) DeleteApplication(applicationName model.ApplicationName) error {
	logger.Debugf("delete application %v", applicationName)
	err := s.applicationDirectory.Delete(applicationName)
	if err != nil {
		logger.Debugf("delete application %v failed: %v", applicationName, err)
		return err
	}
	logger.Debugf("deleted application %v successful", applicationName)
	return nil
}

func (s *service) ExistsApplication(applicationName model.ApplicationName) (bool, error) {
	logger.Debugf("exists application %v", applicationName)
	return s.applicationDirectory.Exists(applicationName)
}

func (s *service) CreateApplication(applicationName model.ApplicationName) (*model.Application, error) {
	logger.Debugf("create application %v", applicationName)
	applicationPassword := model.ApplicationPassword(s.generatePassword(PASSWORD_LENGTH))
	return s.createApplicationWithPassword(applicationName, applicationPassword)
}

func (s *service) CreateApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error) {
	logger.Debugf("create application with password %v", applicationName)
	return s.createApplicationWithPassword(applicationName, applicationPassword)
}

func (s *service) createApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error) {
	exists, err := s.ExistsApplication(applicationName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("applicaton %v already exists", applicationName)
	}
	application := model.Application{
		ApplicationName:     applicationName,
		ApplicationPassword: applicationPassword,
	}
	if err := s.applicationDirectory.Create(applicationName, applicationPassword); err != nil {
		logger.Debugf("create application %v failed: %v", applicationName, err)
		return nil, err
	}
	logger.Debugf("created application %v successful", applicationName)
	return &application, nil
}

func (s *service) VerifyApplicationPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (bool, error) {
	logger.Debugf("verify password of application %v", applicationName)
	pw, err := s.applicationDirectory.Get(applicationName)
	if err != nil {
		return false, err
	}
	if *pw != applicationPassword {
		return false, nil
	}
	return true, nil
}

func (s *service) GetApplication(applicationName model.ApplicationName) (*model.Application, error) {
	logger.Debugf("get application %v", applicationName)
	applicationPassword, err := s.applicationDirectory.Get(applicationName)
	if err != nil {
		return nil, err
	}
	application := model.Application{
		ApplicationName:     applicationName,
		ApplicationPassword: *applicationPassword,
	}
	return &application, nil
}
