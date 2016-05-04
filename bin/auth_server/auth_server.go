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
	"github.com/bborbe/auth/check"
	"github.com/bborbe/auth/filter"
	"github.com/bborbe/auth/ledis"
	"github.com/bborbe/auth/login"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/user_directory"
	flag "github.com/bborbe/flagenv"
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
	PARAMETER_LEDISDB_ADDR              = "ledisdb-address"
	PARAMETER_LEDISDB_PASSWORD          = "ledisdb-password"
)

var (
	logLevelPtr                = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	portPtr                    = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
	ledisdbAddressPtr          = flag.String(PARAMETER_LEDISDB_ADDR, "", "ledisdb address")
	ledisdbPasswordPtr         = flag.String(PARAMETER_LEDISDB_PASSWORD, "", "ledisdb password")
)

func main() {
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	server, err := createServer(*portPtr, *authApplicationPasswordPtr, *ledisdbAddressPtr, *ledisdbPasswordPtr)

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
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
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_LEDISDB_ADDR)
	}
	if len(ledisdbPassword) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_LEDISDB_PASSWORD)
	}

	ledisClient := ledis.New(ledisdbAddress, ledisdbPassword)

	userDirectory := user_directory.New()
	applicationDirectory := application_directory.New(ledisClient)
	applicationCheck := application_check.New(applicationDirectory.Check)

	checkHandler := check.New(ledisClient.Ping)
	accessDeniedHandler := access_denied.New()
	passwordGenerator := generator.New()
	loginHandler := filter.New(applicationCheck.Check, login.New(userDirectory, applicationDirectory.Check).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationCreatorHandler := filter.New(applicationCheck.Check, application_creator.New(applicationDirectory.Create, passwordGenerator.GeneratePassword).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationDeletorHandler := filter.New(applicationCheck.Check, application_deletor.New(applicationDirectory.Delete).ServeHTTP, accessDeniedHandler.ServeHTTP)

	go func() {
		err := applicationDirectory.Create(api.Application{
			ApplicationName:     application_directory.AUTH_APPLICATION_NAME,
			ApplicationPassword: api.ApplicationPassword(authApplicationPassword),
		})
		if err != nil {
			logger.Warnf("create auth application failed: %v", err)
		}
	}()

	handler := router.New(checkHandler.ServeHTTP, loginHandler.ServeHTTP, applicationCreatorHandler.ServeHTTP, applicationDeletorHandler.ServeHTTP)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
