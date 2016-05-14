package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/application_check"
	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/auth/directory/group_user_directory"
	"github.com/bborbe/auth/directory/token_user_directory"
	"github.com/bborbe/auth/directory/user_data_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
	"github.com/bborbe/auth/directory/user_token_directory"
	"github.com/bborbe/auth/filter"
	"github.com/bborbe/auth/handler/access_denied"
	"github.com/bborbe/auth/handler/application_creator"
	"github.com/bborbe/auth/handler/application_deletor"
	"github.com/bborbe/auth/handler/application_getter"
	"github.com/bborbe/auth/handler/check"
	"github.com/bborbe/auth/handler/login"
	"github.com/bborbe/auth/handler/token_adder"
	"github.com/bborbe/auth/handler/token_remover"
	"github.com/bborbe/auth/handler/user_data_delete"
	"github.com/bborbe/auth/handler/user_data_delete_value"
	"github.com/bborbe/auth/handler/user_data_get"
	"github.com/bborbe/auth/handler/user_data_get_value"
	"github.com/bborbe/auth/handler/user_data_set"
	"github.com/bborbe/auth/handler/user_data_set_value"
	"github.com/bborbe/auth/handler/user_delete"
	"github.com/bborbe/auth/handler/user_group_adder"
	"github.com/bborbe/auth/handler/user_group_remover"
	"github.com/bborbe/auth/handler/user_register"
	"github.com/bborbe/auth/handler/user_unregister"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/service/application"
	"github.com/bborbe/auth/service/user"
	"github.com/bborbe/auth/service/user_group"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/ledis"
	"github.com/bborbe/log"
	"github.com/bborbe/password/generator"
	"github.com/facebookgo/grace/gracehttp"
)

var logger = log.DefaultLogger

const (
	DEFAULT_PORT                        = 8080
	PARAMETER_LOGLEVEL                  = "loglevel"
	PARAMETER_PORT                      = "port"
	PARAMETER_AUTH_APPLICATION_PASSWORD = "auth-application-password"
	PARAMETER_LEDISDB_ADDRESS           = "ledisdb-address"
	PARAMETER_LEDISDB_PASSWORD          = "ledisdb-password"
)

var (
	logLevelPtr                = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	portPtr                    = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
	ledisdbAddressPtr          = flag.String(PARAMETER_LEDISDB_ADDRESS, "", "ledisdb address")
	ledisdbPasswordPtr         = flag.String(PARAMETER_LEDISDB_PASSWORD, "", "ledisdb password")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	server, err := createServer(*portPtr, *authApplicationPasswordPtr, *ledisdbAddressPtr, *ledisdbPasswordPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debugf("start server")
	gracehttp.Serve(server)
}

func createServer(port int, authApplicationPassword string, ledisdbAddress string, ledisdbPassword string) (*http.Server, error) {
	logger.Debugf("create server with port: %d", port)
	if port <= 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_PORT)
	}
	if len(ledisdbAddress) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_LEDISDB_ADDRESS)
	}

	ledisClient := ledis.New(ledisdbAddress, ledisdbPassword)
	passwordGenerator := generator.New()

	tokenUserDirectory := token_user_directory.New(ledisClient)
	userTokenDirectory := user_token_directory.New(ledisClient)
	applicationDirectory := application_directory.New(ledisClient)
	groupUserDirectory := group_user_directory.New(ledisClient)
	userGroupDirectory := user_group_directory.New(ledisClient)
	userDataDirectory := user_data_directory.New(ledisClient)

	userService := user.New(userTokenDirectory, userGroupDirectory, tokenUserDirectory)
	applicationService := application.New(passwordGenerator.GeneratePassword, applicationDirectory)
	userGroupService := user_group.New(userGroupDirectory, groupUserDirectory)
	applicationCheck := application_check.New(applicationService.VerifyApplicationPassword)

	checkHandler := check.New(ledisClient.Ping)

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

	userDataSet := user_data_set.New(userDataDirectory.Set)
	userDataSetHandler := filter.New(applicationCheck.Check, userDataSet.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataSetValue := user_data_set_value.New(userDataDirectory.SetValue)
	userDataSetValueHandler := filter.New(applicationCheck.Check, userDataSetValue.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataGet := user_data_get.New(userDataDirectory.Get)
	userDataGetHandler := filter.New(applicationCheck.Check, userDataGet.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataGetValue := user_data_get_value.New(userDataDirectory.GetValue)
	userDataGetValueHandler := filter.New(applicationCheck.Check, userDataGetValue.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataDelete := user_data_delete.New(userDataDirectory.Delete)
	userDataDeleteHandler := filter.New(applicationCheck.Check, userDataDelete.ServeHTTP, accessDeniedHandler.ServeHTTP)

	userDataDeleteValue := user_data_delete_value.New(userDataDirectory.DeleteValue)
	userDataDeleteValueHandler := filter.New(applicationCheck.Check, userDataDeleteValue.ServeHTTP, accessDeniedHandler.ServeHTTP)

	go func() {
		if _, err := applicationService.CreateApplicationWithPassword(api.AUTH_APPLICATION_NAME, api.ApplicationPassword(authApplicationPassword)); err != nil {
			logger.Warnf("create auth application failed: %v", err)
			return
		}
	}()

	handler := router.New(
		checkHandler.ServeHTTP,
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
	)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
