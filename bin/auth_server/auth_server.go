package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bborbe/auth/access_denied"
	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/application_check"
	"github.com/bborbe/auth/application_creator"
	"github.com/bborbe/auth/application_deletor"
	"github.com/bborbe/auth/application_directory"
	"github.com/bborbe/auth/application_getter"
	"github.com/bborbe/auth/application_group_user_directory"
	"github.com/bborbe/auth/application_user_directory"
	"github.com/bborbe/auth/check"
	"github.com/bborbe/auth/filter"
	"github.com/bborbe/auth/login"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/token_adder"
	"github.com/bborbe/auth/token_remover"
	"github.com/bborbe/auth/token_user_directory"
	"github.com/bborbe/auth/user_creator"
	"github.com/bborbe/auth/user_token_directory"
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
	applicationCheck := application_check.New(applicationDirectory.Check)
	applicationUserDirectory := application_user_directory.New()
	applicationGroupUserDirectory := application_group_user_directory.New()
	passwordGenerator := generator.New()
	userCreator := user_creator.New(userTokenDirectory.Add, tokenUserDirectory.Add, userTokenDirectory.Exists, tokenUserDirectory.Exists)
	tokenAdder := token_adder.New(userTokenDirectory.Add, tokenUserDirectory.Add, tokenUserDirectory.Exists, tokenUserDirectory.FindUserByAuthToken)
	tokenRemover := token_remover.New(userTokenDirectory.Remove, tokenUserDirectory.Remove, tokenUserDirectory.FindUserByAuthToken)

	checkHandler := check.New(ledisClient.Ping)
	accessDeniedHandler := access_denied.New()
	loginHandler := filter.New(applicationCheck.Check, login.New(applicationDirectory.Check, tokenUserDirectory.FindUserByAuthToken, tokenUserDirectory.IsUserNotFound, applicationUserDirectory.Contains, applicationGroupUserDirectory.Contains).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationCreatorHandler := filter.New(applicationCheck.Check, application_creator.New(applicationDirectory.Create, passwordGenerator.GeneratePassword).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationDeletorHandler := filter.New(applicationCheck.Check, application_deletor.New(applicationDirectory.Delete).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationGetterHandler := filter.New(applicationCheck.Check, application_getter.New(applicationDirectory.Get, applicationDirectory.IsApplicationNotFound).ServeHTTP, accessDeniedHandler.ServeHTTP)
	userCreateHandler := filter.New(applicationCheck.Check, userCreator.ServeHTTP, accessDeniedHandler.ServeHTTP)
	tokenAddHandler := filter.New(applicationCheck.Check, tokenAdder.ServeHTTP, accessDeniedHandler.ServeHTTP)
	tokenRemoveHandler := filter.New(applicationCheck.Check, tokenRemover.ServeHTTP, accessDeniedHandler.ServeHTTP)

	go func() {
		err := applicationDirectory.Create(api.Application{
			ApplicationName:     application_directory.AUTH_APPLICATION_NAME,
			ApplicationPassword: api.ApplicationPassword(authApplicationPassword),
		})
		if err != nil {
			logger.Warnf("create auth application failed: %v", err)
		}
	}()

	handler := router.New(
		checkHandler.ServeHTTP,
		loginHandler.ServeHTTP,
		applicationCreatorHandler.ServeHTTP,
		applicationDeletorHandler.ServeHTTP,
		applicationGetterHandler.ServeHTTP,
		userCreateHandler.ServeHTTP,
		tokenAddHandler.ServeHTTP,
		tokenRemoveHandler.ServeHTTP,
	)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
