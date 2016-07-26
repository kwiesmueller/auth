package application_check

import (
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/http/header"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type VerifyApplicationPassword func(applicationName model.ApplicationName, applicationPassword model.ApplicationPassword) (bool, error)

type check struct {
	verifyApplicationPassword VerifyApplicationPassword
}

func New(verifyApplicationPassword VerifyApplicationPassword) *check {
	c := new(check)
	c.verifyApplicationPassword = verifyApplicationPassword
	return c
}

func (c *check) Check(req *http.Request) (bool, error) {
	logger.Debugf("validate application")
	name, pass, err := header.ParseAuthorizationBearerHttpRequest(req)
	if err != nil {
		return false, err
	}
	return c.verifyApplicationPassword(model.ApplicationName(name), model.ApplicationPassword(pass))
}
