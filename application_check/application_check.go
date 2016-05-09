package application_check

import (
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/http/bearer"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type VerifyApplicationPassword func(applicationName api.ApplicationName, applicationPassword api.ApplicationPassword) (bool, error)

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
	name, pass, err := bearer.ParseBearerHttpRequest(req)
	if err != nil {
		return false, err
	}
	return c.verifyApplicationPassword(api.ApplicationName(name), api.ApplicationPassword(pass))
}
