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

	"github.com/bborbe/auth/client"
	"github.com/bborbe/auth/factory"
	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/service"
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
	AuthService() service.AuthService
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
		glog.Infof("test normal started")
		fn(createServices(t, config))
		glog.Infof("test normal finished")
	}
	// run rest service
	{
		glog.Infof("test rest started")
		fn(createRestServices(t, config))
		glog.Infof("test rest finished")
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
	return client.New(func(req *http.Request) (*http.Response, error) {
		if req.Body == nil {
			req.Body = ioutil.NopCloser(bytes.NewBufferString(""))
		}
		req.RequestURI = req.URL.RequestURI()
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		return resp.Result(), nil
	}, model.Url("http://example.com"+config.HttpPrefix.String()), config.ApplicationName, config.ApplicationPassword)
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

			username, err := services.AuthService().VerifyTokenHasGroups("testtoken", nil)
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

			username, err := services.AuthService().VerifyTokenHasGroups("testtoken", nil)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username.String(), Is("testuser")); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestHasGroups(t *testing.T) {
	run(t, func(services Services) {
		{

			result, err := services.AuthService().HasGroups("testtoken", nil)
			if err := AssertThat(err, NotNilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(result, Is(false)); err != nil {
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

			result, err := services.AuthService().HasGroups("testtoken", nil)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(result, Is(true)); err != nil {
				t.Fatal(err)
			}

		}
	})
}

func TestHasGroupsWithGroup(t *testing.T) {
	run(t, func(services Services) {
		{
			result, err := services.AuthService().HasGroups("testtoken", []model.GroupName{"testgroup"})
			if err := AssertThat(err, NotNilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(result, Is(false)); err != nil {
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
			err := services.UserGroupService().AddUserToGroup("testuser", "testgroup")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{

			result, err := services.AuthService().HasGroups("testtoken", []model.GroupName{"testgroup"})
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(result, Is(true)); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestAddTokenToUserWithToken(t *testing.T) {
	run(t, func(services Services) {
		{
			err := services.UserService().CreateUserWithToken("testuser", "token1")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserService().AddTokenToUserWithToken("token2", "token1")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			username, err := services.AuthService().VerifyTokenHasGroups("token1", nil)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username.String(), Is("testuser")); err != nil {
				t.Fatal(err)
			}
		}
		{
			username, err := services.AuthService().VerifyTokenHasGroups("token2", nil)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username.String(), Is("testuser")); err != nil {
				t.Fatal(err)
			}
		}
		{
			username, err := services.AuthService().VerifyTokenHasGroups("token3", nil)
			if err := AssertThat(err, NotNilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestDeleteUserWithToken(t *testing.T) {
	run(t, func(services Services) {
		{
			err := services.UserService().CreateUserWithToken("testuser", "token1")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserService().AddTokenToUserWithToken("token2", "token1")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserService().DeleteUserWithToken("token2")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			username, err := services.AuthService().VerifyTokenHasGroups("token1", nil)
			if err := AssertThat(err, NotNilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			username, err := services.AuthService().VerifyTokenHasGroups("token2", nil)
			if err := AssertThat(err, NotNilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestDeleteUserWithTokenSecond(t *testing.T) {
	run(t, func(services Services) {
		{
			err := services.UserService().CreateUserWithToken("testuser", "token1")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserService().DeleteUserWithToken("token1")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			username, err := services.AuthService().VerifyTokenHasGroups("token1", nil)
			if err := AssertThat(err, NotNilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(username, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestListGroupNamesForUsername(t *testing.T) {
	run(t, func(services Services) {
		{
			err := services.UserService().CreateUserWithToken("testuser", "token1")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			groupNames, err := services.UserGroupService().ListGroupNamesForUsername("testuser")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(len(groupNames), Is(0)); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserGroupService().AddUserToGroup("testuser", "testgroup")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			groupNames, err := services.UserGroupService().ListGroupNamesForUsername("testuser")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(len(groupNames), Is(1)); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(groupNames[0].String(), Is("testgroup")); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestAddAndRemoveTokenToUsere(t *testing.T) {
	run(t, func(services Services) {
		token := model.AuthToken("token123")
		username := model.UserName("user123")
		{
			list, err := services.UserService().ListTokenOfUser(username)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(len(list), Is(0)); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserService().AddTokenToUser(token, username)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			list, err := services.UserService().ListTokenOfUser(username)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(len(list), Is(1)); err != nil {
				t.Fatal(err)
			}
		}
		{
			err := services.UserService().RemoveTokenFromUser(token, username)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
		}
		{
			list, err := services.UserService().ListTokenOfUser(username)
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(len(list), Is(0)); err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestCreateApplication(t *testing.T) {
	run(t, func(services Services) {
		{
			app, err := services.ApplicationService().CreateApplication("testapp")
			if err := AssertThat(err, NilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(app, NotNilValue()); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(app.ApplicationName.String(), Is("testapp")); err != nil {
				t.Fatal(err)
			}
			if err := AssertThat(len(app.ApplicationPassword), Gt(0)); err != nil {
				t.Fatal(err)
			}
		}
	})
}
