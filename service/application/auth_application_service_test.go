package application

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/service"
	ledis "github.com/bborbe/redis_client/mock"
)

func TestImplementsService(t *testing.T) {
	object := New(nil, nil)
	var expected *service.ApplicationService
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationCreateOnlyOnce(t *testing.T) {
	applicationDirectory := application_directory.New(ledis.New())
	applicationService := New(func(length int) string {
		return "a"
	}, applicationDirectory)
	name := model.ApplicationName("app")
	application, err := applicationService.CreateApplication(name)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(application, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	application, err = applicationService.CreateApplication(name)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(application, NilValue()); err != nil {
		t.Fatal(err)
	}
}
