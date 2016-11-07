package application

import (
	"github.com/bborbe/auth/model"
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
	panic("not implemented")
}

func (s *applicationService) ExistsApplication(applicationName model.ApplicationName) (bool, error) {
	panic("not implemented")
}

func (s *applicationService) CreateApplication(applicationName model.ApplicationName) (*model.Application, error) {
	panic("not implemented")
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
