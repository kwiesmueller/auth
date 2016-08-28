package application

import (
	"fmt"

	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

const PASSWORD_LENGTH = 16

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
	glog.V(2).Infof("delete application %v", applicationName)
	err := s.applicationDirectory.Delete(applicationName)
	if err != nil {
		glog.V(2).Infof("delete application %v failed: %v", applicationName, err)
		return err
	}
	glog.V(2).Infof("deleted application %v successful", applicationName)
	return nil
}

func (s *service) ExistsApplication(applicationName model.ApplicationName) (bool, error) {
	glog.V(2).Infof("exists application %v", applicationName)
	return s.applicationDirectory.Exists(applicationName)
}

func (s *service) CreateApplication(applicationName model.ApplicationName) (*model.Application, error) {
	glog.V(2).Infof("create application %v", applicationName)
	applicationPassword := model.ApplicationPassword(s.generatePassword(PASSWORD_LENGTH))
	return s.createApplicationWithPassword(applicationName, applicationPassword)
}

func (s *service) CreateApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error) {
	glog.V(2).Infof("create application with password %v", applicationName)
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
		glog.V(2).Infof("create application %v failed: %v", applicationName, err)
		return nil, err
	}
	glog.V(2).Infof("created application %v successful", applicationName)
	return &application, nil
}

func (s *service) VerifyApplicationPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (bool, error) {
	glog.V(2).Infof("verify password of application %v", applicationName)
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
	glog.V(2).Infof("get application %v", applicationName)
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
