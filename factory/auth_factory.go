package factory

import (
	"net/http"

	"github.com/bborbe/auth/application_check"
	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/auth/directory/group_user_directory"
	"github.com/bborbe/auth/directory/token_user_directory"
	"github.com/bborbe/auth/directory/user_data_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
	"github.com/bborbe/auth/directory/user_token_directory"
	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/service"
	"github.com/bborbe/auth/service/application"
	"github.com/bborbe/auth/service/user"
	"github.com/bborbe/auth/service/user_data"
	"github.com/bborbe/auth/service/user_group"
	"github.com/bborbe/auth/v1/handler/access_denied"
	"github.com/bborbe/auth/v1/handler/application_creator"
	"github.com/bborbe/auth/v1/handler/application_deletor"
	"github.com/bborbe/auth/v1/handler/application_getter"
	"github.com/bborbe/auth/v1/handler/login"
	"github.com/bborbe/auth/v1/handler/token_adder"
	"github.com/bborbe/auth/v1/handler/token_remover"
	"github.com/bborbe/auth/v1/handler/user_data_delete"
	"github.com/bborbe/auth/v1/handler/user_data_delete_value"
	"github.com/bborbe/auth/v1/handler/user_data_get"
	"github.com/bborbe/auth/v1/handler/user_data_get_value"
	"github.com/bborbe/auth/v1/handler/user_data_set"
	"github.com/bborbe/auth/v1/handler/user_data_set_value"
	"github.com/bborbe/auth/v1/handler/user_delete"
	"github.com/bborbe/auth/v1/handler/user_group_adder"
	"github.com/bborbe/auth/v1/handler/user_group_remover"
	"github.com/bborbe/auth/v1/handler/user_list"
	"github.com/bborbe/auth/v1/handler/user_register"
	"github.com/bborbe/auth/v1/handler/user_unregister"
	"github.com/bborbe/auth/v1/handler/version"
	"github.com/bborbe/http_handler/check"
	debug_handler "github.com/bborbe/http_handler/debug"
	"github.com/bborbe/http_handler/filter"
	"github.com/bborbe/http_handler/not_found"
	"github.com/bborbe/password/generator"
	"github.com/bborbe/redis_client"
	"github.com/golang/glog"
)

type factory struct {
	config      model.Config
	ledisClient redis_client.Client
}

func New(
	config model.Config,
	ledisClient redis_client.Client,
) *factory {
	h := new(factory)
	h.config = config
	h.ledisClient = ledisClient
	return h
}

func (f *factory) Prefix() model.Prefix { return f.config.Prefix }
func (f *factory) passwordGenerator() generator.PasswordGenerator {
	return generator.New()
}

func (f *factory) tokenUserDirectory() token_user_directory.TokenUserDirectory {
	return token_user_directory.New(f.ledisClient)
}

func (f *factory) userTokenDirectory() user_token_directory.UserTokenDirectory {
	return user_token_directory.New(f.ledisClient)
}

func (f *factory) applicationDirectory() application_directory.ApplicationDirectory {
	return application_directory.New(f.ledisClient)
}

func (f *factory) groupUserDirectory() group_user_directory.GroupUserDirectory {
	return group_user_directory.New(f.ledisClient)
}

func (f *factory) userGroupDirectory() user_group_directory.UserGroupDirectory {
	return user_group_directory.New(f.ledisClient)
}

func (f *factory) userDataDirectory() user_data_directory.UserDataDirectory {
	return user_data_directory.New(f.ledisClient)
}

func (f *factory) ApplicationService() service.ApplicationService {
	return application.New(f.passwordGenerator().GeneratePassword, f.applicationDirectory())
}

func (f *factory) userService() service.UserService {
	return user.New(f.userTokenDirectory(), f.userGroupDirectory(), f.tokenUserDirectory(), f.userDataDirectory())
}

func (f *factory) userGroupService() service.UserGroupService {
	return user_group.New(f.userGroupDirectory(), f.groupUserDirectory())
}

func (f *factory) userDataService() service.UserDataService {
	return user_data.New(f.userDataDirectory())
}

func (f *factory) HealthzHandler() http.Handler {
	return check.New(f.ledisClient.Ping)
}

func (f *factory) ReadinessHandler() http.Handler {
	return f.HealthzHandler()
}

func (f *factory) NotFoundHandler() http.Handler {
	return not_found.New()
}

func (f *factory) Handler() http.Handler {
	handler := router.Create(f)
	if glog.V(4) {
		handler = debug_handler.New(handler)
	}
	return handler
}

func (f *factory) HttpServer() *http.Server {
	glog.V(2).Infof("create http server on %s", f.config.Port.Address())
	return &http.Server{Addr: f.config.Port.Address(), Handler: f.Handler()}
}

func (f *factory) applicationCheck() filter.Check {
	return application_check.New(f.ApplicationService().VerifyApplicationPassword).Check
}

func (f *factory) accessDeniedHandler() http.Handler {
	return access_denied.New()
}

func (f *factory) VersionHandler() http.Handler {
	return version.New()
}

func (f *factory) UserListHandler() http.Handler {
	userList := user_list.New(f.userService().List)
	return filter.New(f.applicationCheck(), userList.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserRegisterHandler() http.Handler {
	userRegister := user_register.New(f.userService().CreateUserWithToken)
	return filter.New(f.applicationCheck(), userRegister.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserDeleteHandler() http.Handler {
	userDelete := user_delete.New(f.userService().DeleteUser)
	return filter.New(f.applicationCheck(), userDelete.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserDataSetHandler() http.Handler {
	userDataSet := user_data_set.New(f.userDataService().Set)
	return filter.New(f.applicationCheck(), userDataSet.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserDataSetValueHandler() http.Handler {
	userDataSetValue := user_data_set_value.New(f.userDataService().SetValue)
	return filter.New(f.applicationCheck(), userDataSetValue.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserDataGetHandler() http.Handler {
	userDataGet := user_data_get.New(f.userDataService().Get)
	return filter.New(f.applicationCheck(), userDataGet.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserDataGetValueHandler() http.Handler {
	userDataGetValue := user_data_get_value.New(f.userDataService().GetValue)
	return filter.New(f.applicationCheck(), userDataGetValue.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserDataDeleteHandler() http.Handler {
	userDataDelete := user_data_delete.New(f.userDataService().Delete)
	return filter.New(f.applicationCheck(), userDataDelete.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserDataDeleteValueHandler() http.Handler {
	userDataDeleteValue := user_data_delete_value.New(f.userDataService().DeleteValue)
	return filter.New(f.applicationCheck(), userDataDeleteValue.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) LoginHandler() http.Handler {
	return filter.New(f.applicationCheck(), login.New(f.userService().VerifyTokenHasGroups).ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) ApplicationCreateHandler() http.Handler {
	return filter.New(f.applicationCheck(), application_creator.New(f.ApplicationService().CreateApplication).ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) ApplicationDeleteHandler() http.Handler {
	return filter.New(f.applicationCheck(), application_deletor.New(f.ApplicationService().DeleteApplication).ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) ApplicationGetHandler() http.Handler {
	return filter.New(f.applicationCheck(), application_getter.New(f.ApplicationService().GetApplication).ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserUnregisterHandler() http.Handler {
	userUnregister := user_unregister.New(f.userService().DeleteUserWithToken)
	return filter.New(f.applicationCheck(), userUnregister.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) TokenAddHandler() http.Handler {
	tokenAdder := token_adder.New(f.userService().AddTokenToUserWithToken)
	return filter.New(f.applicationCheck(), tokenAdder.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) TokenRemoveHandler() http.Handler {
	tokenRemover := token_remover.New(f.userService().RemoveTokenFromUserWithToken)
	return filter.New(f.applicationCheck(), tokenRemover.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserGroupAddHandler() http.Handler {
	userGroupAdder := user_group_adder.New(f.userGroupService().AddUserToGroup)
	return filter.New(f.applicationCheck(), userGroupAdder.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) UserGroupRemoveHandler() http.Handler {
	userGroupRemover := user_group_remover.New(f.userGroupService().RemoveUserFromGroup)
	return filter.New(f.applicationCheck(), userGroupRemover.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}
