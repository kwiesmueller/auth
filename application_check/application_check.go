package application_check

import (
	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/http/header"
	"github.com/golang/glog"
)

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
	glog.V(3).Infof("validate application started")
	name, pass, err := header.ParseAuthorizationBearerHttpRequest(req)
	if err != nil {
		glog.V(2).Infof("parse header failed: %v", err)
		return false, err
	}
	result, err := c.verifyApplicationPassword(model.ApplicationName(name), model.ApplicationPassword(pass))
	if err != nil {
		glog.V(2).Infof("verify application password failed: %v", err)
		return false, err
	}
	glog.V(3).Infof("validate application finished")
	return result, nil
}
