package main

import (
	"testing"

	. "github.com/bborbe/assert"
)
import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/bborbe/auth/client/application"
	"github.com/bborbe/auth/client/rest"
	"github.com/bborbe/auth/client/user"
	"github.com/bborbe/auth/client/user_data"
	"github.com/bborbe/auth/client/user_group"
	"github.com/bborbe/auth/factory"
	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/service"
	http_rest "github.com/bborbe/http/rest"
	ledis_mock "github.com/bborbe/redis_client/mock"
	"github.com/golang/glog"
)

func TestServerConfig(t *testing.T) {
	config := createConfig()
	if err := AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

type Services interface {
	ApplicationService() service.ApplicationService
	UserService() service.UserService
	UserDataService() service.UserDataService
	UserGroupService() service.UserGroupService
}

func run(t *testing.T, fn func(services Services)) {
	config := model.Config{
		HttpPrefix:          "/auth",
		ApplicationName:     "testauth",
		ApplicationPassword: "test123",
	}
	// run normal service
	{
		fn(createServices(t, config))
	}
	// run rest service
	{
		fn(createRestServices(t, config))
	}
}

type restServices struct {
	applicationService service.ApplicationService
	userService        service.UserService
	userDataService    service.UserDataService
	userGroupService   service.UserGroupService
}

func createServices(t *testing.T, config model.Config) Services {
	factory := factory.New(config, ledis_mock.New())
	if _, err := factory.ApplicationService().CreateApplicationWithPassword(config.ApplicationName, config.ApplicationPassword); err != nil {
		t.Fatal("create app failed", err)
	}
	return factory
}

func createRestServices(t *testing.T, config model.Config) Services {
	factory := factory.New(config, ledis_mock.New())
	if _, err := factory.ApplicationService().CreateApplicationWithPassword(config.ApplicationName, config.ApplicationPassword); err != nil {
		t.Fatal("create app failed", err)
	}

	handler := router.Create(factory)
	httpRest := http_rest.New(func(req *http.Request) (*http.Response, error) {
		if req.Body == nil {
			req.Body = ioutil.NopCloser(bytes.NewBufferString(""))
		}
		req.RequestURI = req.URL.RequestURI()
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		return resp.Result(), nil
	})
	rest := rest.New(httpRest.Call, model.Url("http://example.com"+config.HttpPrefix.String()), config.ApplicationName, config.ApplicationPassword)
	return &restServices{
		applicationService: application.New(rest.Call),
		userService:        user.New(rest.Call),
		userDataService:    user_data.New(rest.Call),
		userGroupService:   user_group.New(rest.Call),
	}
}

func (r *restServices) ApplicationService() service.ApplicationService {
	return r.applicationService
}
func (r *restServices) UserService() service.UserService {
	return r.userService
}
func (r *restServices) UserDataService() service.UserDataService {
	return r.userDataService
}
func (r *restServices) UserGroupService() service.UserGroupService {
	return r.userGroupService
}

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestUserServiceList(t *testing.T) {
	run(t, func(services Services) {
		{
			list, err := services.UserService().List()
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(len(list), Is(0)); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserService().CreateUserWithToken("testuser", "testtoken")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			list, err := services.UserService().List()
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(len(list), Is(1)); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(list[0].String(), Is("testuser")); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestVerifyTokenHasGroups(t *testing.T) {
	run(t, func(services Services) {
		{

			username, err := services.UserService().VerifyTokenHasGroups("testtoken", nil)
			if err := AssertThat(err, NotNilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserService().CreateUserWithToken("testuser", "testtoken")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{

			username, err := services.UserService().VerifyTokenHasGroups("testtoken", nil)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username.String(), Is("testuser")); err != nil {
				t.Fatal(err)
			}
		}
	})
}
