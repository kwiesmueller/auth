package router

import (
	"testing"

	. "github.com/bborbe/assert"

	"net/http"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/http/mock"
	"github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/http_handler/adapter"
)

type counter struct {
	notFound              int
	healthz               int
	readiness             int
	login                 int
	applicationCreate     int
	applicationDelete     int
	applicationGet        int
	userRegister          int
	userUnregister        int
	userDelete            int
	tokenAdd              int
	tokenRemove           int
	userGroupAdd          int
	userGroupRemove       int
	userDataSet           int
	userDataSetValue      int
	userDataGet           int
	userDataGetValue      int
	userDataDelete        int
	userDataDeleteValue   int
	version               int
	userList              int
	tokensForUsername     int
	groupNamesForUsername int
}

func createCounterHandler(counter *int) http.Handler {
	return adapter.New(func(http.ResponseWriter, *http.Request) {
		*counter++
	})
}

func (c *counter) Prefix() model.Prefix {
	return "/prefix"
}

func (c *counter) TokensForUsernameHandler() http.Handler {
	return createCounterHandler(&c.tokensForUsername)
}

func (c *counter) NotFoundHandler() http.Handler {
	return createCounterHandler(&c.notFound)
}

func (c *counter) HealthzHandler() http.Handler {
	return createCounterHandler(&c.healthz)
}

func (c *counter) ReadinessHandler() http.Handler {
	return createCounterHandler(&c.readiness)
}

func (c *counter) VersionHandler() http.Handler {
	return createCounterHandler(&c.version)
}

func (c *counter) UserListHandler() http.Handler {
	return createCounterHandler(&c.userList)
}

func (c *counter) UserRegisterHandler() http.Handler {
	return createCounterHandler(&c.userRegister)
}

func (c *counter) UserDeleteHandler() http.Handler {
	return createCounterHandler(&c.userDelete)
}

func (c *counter) UserDataSetHandler() http.Handler {
	return createCounterHandler(&c.userDataSet)
}

func (c *counter) UserDataSetValueHandler() http.Handler {
	return createCounterHandler(&c.userDataSetValue)
}

func (c *counter) UserDataGetHandler() http.Handler {
	return createCounterHandler(&c.userDataGet)
}

func (c *counter) UserDataGetValueHandler() http.Handler {
	return createCounterHandler(&c.userDataGetValue)
}

func (c *counter) GroupNamesForUsernameHandler() http.Handler {
	return createCounterHandler(&c.groupNamesForUsername)
}

func (c *counter) UserDataDeleteHandler() http.Handler {
	return createCounterHandler(&c.userDataDelete)
}

func (c *counter) UserDataDeleteValueHandler() http.Handler {
	return createCounterHandler(&c.userDataDeleteValue)
}

func (c *counter) LoginHandler() http.Handler {
	return createCounterHandler(&c.login)
}

func (c *counter) ApplicationCreateHandler() http.Handler {
	return createCounterHandler(&c.applicationCreate)
}

func (c *counter) ApplicationDeleteHandler() http.Handler {
	return createCounterHandler(&c.applicationDelete)
}

func (c *counter) ApplicationGetHandler() http.Handler {
	return createCounterHandler(&c.applicationGet)
}

func (c *counter) UserUnregisterHandler() http.Handler {
	return createCounterHandler(&c.userUnregister)
}

func (c *counter) TokenAddHandler() http.Handler {
	return createCounterHandler(&c.tokenAdd)
}

func (c *counter) TokenRemoveHandler() http.Handler {
	return createCounterHandler(&c.tokenRemove)
}

func (c *counter) UserGroupAddHandler() http.Handler {
	return createCounterHandler(&c.userGroupAdd)
}

func (c *counter) UserGroupRemoveHandler() http.Handler {
	return createCounterHandler(&c.userGroupRemove)
}

func TestHealthz(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/healthz")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.healthz, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.healthz, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestReadiness(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/readiness")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.readiness, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.readiness, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestNotFound(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/asdf")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.notFound, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.notFound, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/login")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.login, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.login, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationCreate(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/application")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.applicationCreate, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.applicationCreate, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationDelete(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/application/test123")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.applicationDelete, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.applicationDelete, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationGet(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/application/test123")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.applicationGet, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.applicationGet, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserRegister(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userRegister, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userRegister, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserDelete(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user/123")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userDelete, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userDelete, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUnregister(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/token/123")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userUnregister, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userUnregister, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestTokenAdd(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/token")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.tokenAdd, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.tokenAdd, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestTokenRemove(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/token")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.tokenRemove, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.tokenRemove, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserGroupAdd(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user_group")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userGroupAdd, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userGroupAdd, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestGroupNamesForUsername(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user_group?username=foo")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.groupNamesForUsername, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.groupNamesForUsername, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserGroupRemove(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user_group")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userGroupRemove, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userGroupRemove, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserDataSet(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user/tester/data")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userDataSet, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userDataSet, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserDataSetValue(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user/tester/data/keya")
	rb.SetMethod("POST")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userDataSetValue, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userDataSetValue, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserDataGet(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user/tester/data")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userDataGet, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userDataGet, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserDataGetValue(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user/tester/data/keya")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userDataGetValue, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userDataGetValue, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserDataDelete(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user/tester/data")
	rb.SetMethod("DELETE")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userDataDelete, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userDataDelete, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserDataDeleteValue(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()
	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user/tester/data/keya")
	rb.SetMethod("DELETE")

	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userDataDeleteValue, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userDataDeleteValue, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestVersion(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/version")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.version, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.version, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestUserList(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/user")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.userList, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.userList, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestTokensForUsername(t *testing.T) {
	c := new(counter)
	r := Create(c)
	resp := mock.NewHttpResponseWriterMock()

	rb := requestbuilder.NewHTTPRequestBuilder("http://example.com/prefix/api/1.0/token?username=test123")
	rb.SetMethod("GET")
	req, err := rb.Build()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(c.tokensForUsername, Is(0)); err != nil {
		t.Fatal(err)
	}
	r.ServeHTTP(resp, req)
	if err := AssertThat(c.tokensForUsername, Is(1)); err != nil {
		t.Fatal(err)
	}
}
