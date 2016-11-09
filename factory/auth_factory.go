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
	"github.com/bborbe/auth/v1/handler/groupnames_by_username"
	"github.com/bborbe/auth/v1/handler/login"
	"github.com/bborbe/auth/v1/handler/token_adder"
	"github.com/bborbe/auth/v1/handler/token_remover"
	"github.com/bborbe/auth/v1/handler/tokens_by_username"
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

func (f *factory) addRequireAuth(handler http.Handler) http.Handler {
	return filter.New(f.applicationCheck(), handler.ServeHTTP, f.accessDeniedHandler().ServeHTTP)
}

func (f *factory) ApplicationService() service.ApplicationService {
	return application.New(f.passwordGenerator().GeneratePassword, f.applicationDirectory())
}

func (f *factory) UserService() service.UserService {
	return user.New(f.userTokenDirectory(), f.userGroupDirectory(), f.tokenUserDirectory(), f.userDataDirectory())
}

func (f *factory) UserGroupService() service.UserGroupService {
	return user_group.New(f.userGroupDirectory(), f.groupUserDirectory())
}

func (f *factory) UserDataService() service.UserDataService {
	return user_data.New(f.userDataDirectory())
}

func (f *factory) Prefix() model.Prefix {
	return f.config.HttpPrefix
}
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
	glog.V(2).Infof("create http server on %s", f.config.HttpPort.Address())
	return &http.Server{Addr: f.config.HttpPort.Address(), Handler: f.Handler()}
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
	return f.addRequireAuth(user_list.New(f.UserService().List))
}

func (f *factory) UserRegisterHandler() http.Handler {
	return f.addRequireAuth(user_register.New(f.UserService().CreateUserWithToken))
}

func (f *factory) UserDeleteHandler() http.Handler {
	return f.addRequireAuth(user_delete.New(f.UserService().DeleteUser))
}

func (f *factory) UserDataSetHandler() http.Handler {
	return f.addRequireAuth(user_data_set.New(f.UserDataService().Set))
}

func (f *factory) UserDataSetValueHandler() http.Handler {
	return f.addRequireAuth(user_data_set_value.New(f.UserDataService().SetValue))
}

func (f *factory) UserDataGetHandler() http.Handler {
	return f.addRequireAuth(user_data_get.New(f.UserDataService().Get))
}

func (f *factory) UserDataGetValueHandler() http.Handler {
	return f.addRequireAuth(user_data_get_value.New(f.UserDataService().GetValue))
}

func (f *factory) UserDataDeleteHandler() http.Handler {
	return f.addRequireAuth(user_data_delete.New(f.UserDataService().Delete))
}

func (f *factory) UserDataDeleteValueHandler() http.Handler {
	return f.addRequireAuth(user_data_delete_value.New(f.UserDataService().DeleteValue))
}

func (f *factory) LoginHandler() http.Handler {
	return f.addRequireAuth(login.New(f.UserService().VerifyTokenHasGroups))
}

func (f *factory) ApplicationCreateHandler() http.Handler {
	return f.addRequireAuth(application_creator.New(f.ApplicationService().CreateApplication))
}

func (f *factory) ApplicationDeleteHandler() http.Handler {
	return f.addRequireAuth(application_deletor.New(f.ApplicationService().DeleteApplication))
}

func (f *factory) ApplicationGetHandler() http.Handler {
	return f.addRequireAuth(application_getter.New(f.ApplicationService().GetApplication))
}

func (f *factory) UserUnregisterHandler() http.Handler {
	return f.addRequireAuth(user_unregister.New(f.UserService().DeleteUserWithToken))
}

func (f *factory) TokenAddHandler() http.Handler {
	return f.addRequireAuth(token_adder.New(f.UserService().AddTokenToUserWithToken))
}

func (f *factory) TokenRemoveHandler() http.Handler {
	return f.addRequireAuth(token_remover.New(f.UserService().RemoveTokenFromUserWithToken))
}

func (f *factory) UserGroupAddHandler() http.Handler {
	return f.addRequireAuth(user_group_adder.New(f.UserGroupService().AddUserToGroup))
}

func (f *factory) UserGroupRemoveHandler() http.Handler {
	return f.addRequireAuth(user_group_remover.New(f.UserGroupService().RemoveUserFromGroup))
}

func (f *factory) TokensForUsernameHandler() http.Handler {
	return f.addRequireAuth(tokens_by_username.New(f.UserService().ListTokenOfUser))
}

func (f *factory) GroupNamesForUsernameHandler() http.Handler {
	return f.addRequireAuth(groupnames_by_username.New(f.UserGroupService().ListGroupNamesForUsername))
}
