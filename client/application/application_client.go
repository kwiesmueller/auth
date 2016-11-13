package application

import (
	"fmt"

	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/golang/glog"
)

type callRest func(path string, method string, request interface{}, response interface{}) error

type applicationService struct {
	callRest callRest
}

func New(
	callRest callRest,
) *applicationService {
	s := new(applicationService)
	s.callRest = callRest
	return s
}

func (s *applicationService) DeleteApplication(applicationName model.ApplicationName) error {
	glog.V(3).Infof("delete application %s", applicationName)
	if err := s.callRest(fmt.Sprintf("/api/1.0/application/%s", applicationName), http.MethodDelete, nil, nil); err != nil {
		glog.V(2).Infof("delete application %s failed: %v", applicationName, err)
		return err
	}
	glog.V(3).Infof("delete application %s successful", applicationName)
	return nil
}

func (s *applicationService) ExistsApplication(applicationName model.ApplicationName) (bool, error) {
	glog.V(3).Infof("exists application %s", applicationName)
	var response model.Application
	if err := s.callRest(fmt.Sprintf("/api/1.0/application/%s", applicationName), http.MethodGet, nil, &response); err != nil {
		glog.V(2).Infof("exists application %s failed: %v", applicationName, err)
		return false, err
	}
	glog.V(3).Infof("exists application %s successful", applicationName)
	return len(response.ApplicationPassword) > 0, nil
}

func (s *applicationService) CreateApplication(applicationName model.ApplicationName) (*model.Application, error) {
	glog.V(3).Infof("create application %s", applicationName)
	request := model.Application{
		ApplicationName: applicationName,
	}
	var response model.Application
	if err := s.callRest("/api/1.0/application", http.MethodPost, &request, &response); err != nil {
		glog.V(2).Infof("create application %s failed: %v", applicationName, err)
		return nil, err
	}
	glog.V(3).Infof("create application %s successful", applicationName)
	return &response, nil
}

func (s *applicationService) CreateApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error) {
	glog.V(3).Infof("create application with password %s", applicationName)
	request := model.Application{
		ApplicationName:     applicationName,
		ApplicationPassword: applicationPassword,
	}
	var response model.Application
	if err := s.callRest("/api/1.0/application", http.MethodPost, &request, &response); err != nil {
		glog.V(2).Infof("create application with password %s failed: %v", applicationName, err)
		return nil, err
	}
	glog.V(3).Infof("create application with password %s successful", applicationName)
	return &response, nil
}

func (s *applicationService) VerifyApplicationPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (s *applicationService) GetApplication(applicationName model.ApplicationName) (*model.Application, error) {
	glog.V(3).Infof("exists application %s", applicationName)
	var response model.Application
	if err := s.callRest(fmt.Sprintf("/api/1.0/application/%s", applicationName), http.MethodGet, nil, &response); err != nil {
		glog.V(2).Infof("exists application %s failed: %v", applicationName, err)
		return nil, err
	}
	glog.V(3).Infof("exists application %s successful", applicationName)
	return &response, nil
}
