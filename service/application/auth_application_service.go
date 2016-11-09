package application

import (
	"fmt"

	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

const passwordLength = 16

type generatePassword func(length int) string

type applicationService struct {
	generatePassword     generatePassword
	applicationDirectory application_directory.ApplicationDirectory
}

func New(generatePassword generatePassword, applicationDirectory application_directory.ApplicationDirectory) *applicationService {
	s := new(applicationService)
	s.generatePassword = generatePassword
	s.applicationDirectory = applicationDirectory
	return s
}

func (s *applicationService) DeleteApplication(applicationName model.ApplicationName) error {
	glog.V(4).Infof("delete application %v", applicationName)
	err := s.applicationDirectory.Delete(applicationName)
	if err != nil {
		glog.V(2).Infof("delete application %v failed: %v", applicationName, err)
		return err
	}
	glog.V(4).Infof("deleted application %v successful", applicationName)
	return nil
}

func (s *applicationService) ExistsApplication(applicationName model.ApplicationName) (bool, error) {
	glog.V(4).Infof("exists application %v", applicationName)
	return s.applicationDirectory.Exists(applicationName)
}

func (s *applicationService) CreateApplication(applicationName model.ApplicationName) (*model.Application, error) {
	glog.V(4).Infof("create application %v", applicationName)
	applicationPassword := model.ApplicationPassword(s.generatePassword(passwordLength))
	return s.createApplicationWithPassword(applicationName, applicationPassword)
}

func (s *applicationService) CreateApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error) {
	glog.V(4).Infof("create application with password %v", applicationName)
	return s.createApplicationWithPassword(applicationName, applicationPassword)
}

func (s *applicationService) createApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error) {
	glog.V(4).Infof("create application %v with password-length %d", applicationName, len(applicationPassword))
	exists, err := s.ExistsApplication(applicationName)
	if err != nil {
		glog.V(2).Infof("check application exists failed: %v", err)
		return nil, err
	}
	if exists {
		glog.V(2).Infof("applicaton %v already exists", applicationName)
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
	glog.V(4).Infof("created application %v successful", applicationName)
	return &application, nil
}

func (s *applicationService) VerifyApplicationPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (bool, error) {
	glog.V(4).Infof("verify password of application %v", applicationName)
	pw, err := s.applicationDirectory.Get(applicationName)
	if err != nil {
		glog.V(2).Infof("get application %v failed: %v", applicationName, err)
		return false, err
	}
	if *pw != applicationPassword {
		glog.V(2).Infof("password for application %v invalid", applicationName)
		return false, nil
	}
	glog.V(4).Infof("password for application %v valid", applicationName)
	return true, nil
}

func (s *applicationService) GetApplication(applicationName model.ApplicationName) (*model.Application, error) {
	glog.V(4).Infof("get application %v", applicationName)
	applicationPassword, err := s.applicationDirectory.Get(applicationName)
	if err != nil {
		glog.V(2).Infof("get application %v failed: %v", applicationName, err)
		return nil, err
	}
	application := model.Application{
		ApplicationName:     applicationName,
		ApplicationPassword: *applicationPassword,
	}
	return &application, nil
}
