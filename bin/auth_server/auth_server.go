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
	"github.com/bborbe/auth/handler/user_register"
	"github.com/bborbe/auth/handler/user_unregister"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/service/user"
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
	PARAMETER_ADMIN                     = "admin"
	PARAMETER_AUTH_APPLICATION_PASSWORD = "auth-application-password"
	PARAMETER_LEDISDB_ADDRESS           = "ledisdb-address"
	PARAMETER_LEDISDB_PASSWORD          = "ledisdb-password"
)

var (
	logLevelPtr                = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	portPtr                    = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
	adminUserNamePtr           = flag.String(PARAMETER_ADMIN, "", "admin username")
	ledisdbAddressPtr          = flag.String(PARAMETER_LEDISDB_ADDRESS, "", "ledisdb address")
	ledisdbPasswordPtr         = flag.String(PARAMETER_LEDISDB_PASSWORD, "", "ledisdb password")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	server, err := createServer(*portPtr, *authApplicationPasswordPtr, *adminUserNamePtr, *ledisdbAddressPtr, *ledisdbPasswordPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debugf("start server")
	gracehttp.Serve(server)
}

func createServer(port int, authApplicationPassword string, adminUserName string, ledisdbAddress string, ledisdbPassword string) (*http.Server, error) {
	logger.Debugf("create server with port: %d", port)
	if port <= 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_PORT)
	}
	if len(authApplicationPassword) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_AUTH_APPLICATION_PASSWORD)
	}
	if len(ledisdbAddress) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_LEDISDB_ADDRESS)
	}

	ledisClient := ledis.New(ledisdbAddress, ledisdbPassword)

	tokenUserDirectory := token_user_directory.New(ledisClient)
	userTokenDirectory := user_token_directory.New(ledisClient)
	applicationDirectory := application_directory.New(ledisClient)
	groupUserDirectory := group_user_directory.New(ledisClient)
	userGroupDirectory := user_group_directory.New(ledisClient)

	userService := user.New(userTokenDirectory, userGroupDirectory, tokenUserDirectory)

	applicationCheck := application_check.New(applicationDirectory.Check)
	passwordGenerator := generator.New()
	userRegister := user_register.New(userService.CreateUserWithToken)
	userUnregister := user_unregister.New(userService.DeleteUserWithToken)
	tokenAdder := token_adder.New(userService.AddTokenToUserWithToken)
	tokenRemover := token_remover.New(userService.RemoveTokenFromUserWithToken)

	checkHandler := check.New(ledisClient.Ping)
	accessDeniedHandler := access_denied.New()
	loginHandler := filter.New(applicationCheck.Check, login.New(userService.VerifyTokenHasGroups).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationCreatorHandler := filter.New(applicationCheck.Check, application_creator.New(applicationDirectory.Create, passwordGenerator.GeneratePassword).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationDeletorHandler := filter.New(applicationCheck.Check, application_deletor.New(applicationDirectory.Delete).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationGetterHandler := filter.New(applicationCheck.Check, application_getter.New(applicationDirectory.Get, applicationDirectory.IsApplicationNotFound).ServeHTTP, accessDeniedHandler.ServeHTTP)
	userRegisterHandler := filter.New(applicationCheck.Check, userRegister.ServeHTTP, accessDeniedHandler.ServeHTTP)
	userUnregisterHandler := filter.New(applicationCheck.Check, userUnregister.ServeHTTP, accessDeniedHandler.ServeHTTP)
	tokenAddHandler := filter.New(applicationCheck.Check, tokenAdder.ServeHTTP, accessDeniedHandler.ServeHTTP)
	tokenRemoveHandler := filter.New(applicationCheck.Check, tokenRemover.ServeHTTP, accessDeniedHandler.ServeHTTP)

	go func() {
		var err error
		err = applicationDirectory.Create(api.Application{
			ApplicationName:     api.AUTH_APPLICATION_NAME,
			ApplicationPassword: api.ApplicationPassword(authApplicationPassword),
		})
		if err != nil {
			logger.Warnf("create auth application failed: %v", err)
			return
		}
		err = groupUserDirectory.Add(api.AUTH_ADMIN_GROUP, api.UserName(adminUserName))
		if err != nil {
			logger.Warnf("add user %s to group %v failed: %v", adminUserName, api.AUTH_ADMIN_GROUP, err)
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
		tokenAddHandler.ServeHTTP,
		tokenRemoveHandler.ServeHTTP,
	)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
