package handler_creator

import (
	"net/http"

	"github.com/bborbe/auth/application_check"
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
	"github.com/bborbe/auth/v1/router"
	"github.com/bborbe/auth/v1/version"
	"github.com/bborbe/http_handler/filter"
	"github.com/bborbe/http_handler/not_found"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type HandlerCreator interface {
	CreateHandler(
		prefix string,
		applicationService application.Service,
		userService user.Service,
		userGroupService user_group.Service,
		userDataService user_data.UserDataService,
	) (http.Handler, error)
}

type handlerCreator struct {
}

func New() *handlerCreator {
	h := new(handlerCreator)
	return h
}

func (h *handlerCreator) CreateHandler(
	prefix string,
	applicationService application.Service,
	userService user.Service,
	userGroupService user_group.Service,
	userDataService user_data.UserDataService,
) (http.Handler, error) {

	applicationCheck := application_check.New(applicationService.VerifyApplicationPassword)

	accessDeniedHandler := access_denied.New()

	loginHandler := filter.New(applicationCheck.Check, login.New(userService.VerifyTokenHasGroups).ServeHTTP, accessDeniedHandler.ServeHTTP)

	applicationCreatorHandler := filter.New(applicationCheck.Check, application_creator.New(applicationService.CreateApplication).ServeHTTP, accessDeniedHandler.ServeHTTP)

	applicationDeletorHandler := filter.New(applicationCheck.Check, application_deletor.New(applicationService.DeleteApplication).ServeHTTP, accessDeniedHandler.ServeHTTP)

	applicationGetterHandler := filter.New(applicationCheck.Check, application_getter.New(applicationService.GetApplication).ServeHTTP, accessDeniedHandler.ServeHTTP)

	userRegister := user_register.New(userService.CreateUserWithToken)
	userRegisterHandler := filter.New(applicationCheck.Check, userRegister.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userUnregister := user_unregister.New(userService.DeleteUserWithToken)
	userUnregisterHandler := filter.New(applicationCheck.Check, userUnregister.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDelete := user_delete.New(userService.DeleteUser)
	userDeleteHandler := filter.New(applicationCheck.Check, userDelete.ServeHTTP, accessDeniedHandler.ServeHTTP)

	tokenAdder := token_adder.New(userService.AddTokenToUserWithToken)
	tokenAddHandler := filter.New(applicationCheck.Check, tokenAdder.ServeHTTP, accessDeniedHandler.ServeHTTP)

	tokenRemover := token_remover.New(userService.RemoveTokenFromUserWithToken)
	tokenRemoveHandler := filter.New(applicationCheck.Check, tokenRemover.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userGroupAdder := user_group_adder.New(userGroupService.AddUserToGroup)
	userGroupAddHandler := filter.New(applicationCheck.Check, userGroupAdder.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userGroupRemover := user_group_remover.New(userGroupService.RemoveUserFromGroup)
	userGroupRemoveHandler := filter.New(applicationCheck.Check, userGroupRemover.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataSet := user_data_set.New(userDataService.Set)
	userDataSetHandler := filter.New(applicationCheck.Check, userDataSet.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataSetValue := user_data_set_value.New(userDataService.SetValue)
	userDataSetValueHandler := filter.New(applicationCheck.Check, userDataSetValue.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataGet := user_data_get.New(userDataService.Get)
	userDataGetHandler := filter.New(applicationCheck.Check, userDataGet.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataGetValue := user_data_get_value.New(userDataService.GetValue)
	userDataGetValueHandler := filter.New(applicationCheck.Check, userDataGetValue.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataDelete := user_data_delete.New(userDataService.Delete)
	userDataDeleteHandler := filter.New(applicationCheck.Check, userDataDelete.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataDeleteValue := user_data_delete_value.New(userDataService.DeleteValue)
	userDataDeleteValueHandler := filter.New(applicationCheck.Check, userDataDeleteValue.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userList := user_list.New(userService.List)
	userListHandler := filter.New(applicationCheck.Check, userList.ServeHTTP, accessDeniedHandler.ServeHTTP)

	notFoundHandler := not_found.New()

	versionHandler := version.New()

	handler := router.New(
		prefix,
		notFoundHandler.ServeHTTP,
		versionHandler.ServeHTTP,
		loginHandler.ServeHTTP,
		applicationCreatorHandler.ServeHTTP,
		applicationDeletorHandler.ServeHTTP,
		applicationGetterHandler.ServeHTTP,
		userRegisterHandler.ServeHTTP,
		userUnregisterHandler.ServeHTTP,
		userDeleteHandler.ServeHTTP,
		tokenAddHandler.ServeHTTP,
		tokenRemoveHandler.ServeHTTP,
		userGroupAddHandler.ServeHTTP,
		userGroupRemoveHandler.ServeHTTP,
		userDataSetHandler.ServeHTTP,
		userDataSetValueHandler.ServeHTTP,
		userDataGetHandler.ServeHTTP,
		userDataGetValueHandler.ServeHTTP,
		userDataDeleteHandler.ServeHTTP,
		userDataDeleteValueHandler.ServeHTTP,
		userListHandler.ServeHTTP,
	)
	return handler, nil
}
