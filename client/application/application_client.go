package application

import (
	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
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
	glog.V(4).Infof("delete application %s", applicationName)
	if err := s.callRest(fmt.Sprintf("/api/1.0/application/%s", applicationName), "DELETE", nil, nil); err != nil {
		glog.V(2).Infof("delete application %s failed: %v", applicationName, err)
		return err
	}
	glog.V(4).Infof("delete application %s successful", applicationName)
	return nil
}

func (s *applicationService) ExistsApplication(applicationName model.ApplicationName) (bool, error) {
	glog.V(4).Infof("exists application %s", applicationName)
	var response v1.GetApplicationResponse
	if err := s.callRest(fmt.Sprintf("/api/1.0/application/%s", applicationName), "GET", nil, &response); err != nil {
		glog.V(2).Infof("exists application %s failed: %v", applicationName, err)
		return false, err
	}
	glog.V(4).Infof("exists application %s successful", applicationName)
	return len(response.ApplicationPassword) > 0, nil
}

func (s *applicationService) CreateApplication(applicationName model.ApplicationName) (*model.Application, error) {
	glog.V(4).Infof("create application %s", applicationName)
	request := v1.CreateApplicationRequest{
		ApplicationName: model.ApplicationName(applicationName),
	}
	var response model.Application
	if err := s.callRest("/api/1.0/application", "POST", &request, &response); err != nil {
		glog.V(2).Infof("create application %s failed: %v", applicationName, err)
		return nil, err
	}
	glog.V(4).Infof("create application %s successful", applicationName)
	return &response, nil
}

func (s *applicationService) CreateApplicationWithPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (*model.Application, error) {
	panic("not implemented")
}

func (s *applicationService) VerifyApplicationPassword(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (bool, error) {
	panic("not implemented")
}

func (s *applicationService) GetApplication(applicationName model.ApplicationName) (*model.Application, error) {
	panic("not implemented")
}
