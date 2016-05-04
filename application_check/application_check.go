package application_check

import (
	"net/http"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/http/bearer"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type CheckApplication func(api.ApplicationName, api.ApplicationPassword) error

type check struct {
	checkApplication CheckApplication
}

func New(checkApplication CheckApplication) *check {
	c := new(check)
	c.checkApplication = checkApplication
	return c
}

func (c *check) Check(req *http.Request) (bool, error) {
	logger.Debugf("validate application")
	name, pass, err := bearer.ParseBearerHttpRequest(req)
	if err != nil {
		return false, err
	}
	return c.checkApplication(api.ApplicationName(name), api.ApplicationPassword(pass)) == nil, nil
}
