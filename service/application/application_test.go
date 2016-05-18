package application

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/ledis/mock"
)

func TestImplementsService(t *testing.T) {
	object := New(nil, nil)
	var expected *Service
	if err := AssertThat(object, Implements(expected)); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationCreateOnlyOnce(t *testing.T) {
	ledis := mock.NewKv()
	applicationDirectory := application_directory.New(ledis)
	applicationService := New(func(length int) string {
		return "a"
	}, applicationDirectory)
	name := api.ApplicationName("app")
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
