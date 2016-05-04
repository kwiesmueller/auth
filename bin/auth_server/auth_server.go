package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bborbe/auth/access_denied"
	"github.com/bborbe/auth/api"
	"github.com/bborbe/auth/application_check"
	"github.com/bborbe/auth/application_creator"
	"github.com/bborbe/auth/application_directory"
	"github.com/bborbe/auth/check"
	"github.com/bborbe/auth/filter"
	"github.com/bborbe/auth/login"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/user_directory"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/log"
	"github.com/facebookgo/grace/gracehttp"
)

var logger = log.DefaultLogger

const (
	DEFAULT_PORT                        = 8080
	PARAMETER_LOGLEVEL                  = "loglevel"
	PARAMETER_PORT                      = "port"
	PARAMETER_AUTH_APPLICATION_PASSWORD = "auth-application-password"
)

var (
	logLevelPtr                = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	portPtr                    = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
)

func main() {
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	server, err := createServer(*portPtr, *authApplicationPasswordPtr)

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	gracehttp.Serve(server)
}

func createServer(port int, authApplicationPassword string) (*http.Server, error) {
	logger.Debugf("create server with port: %d", port)
	if port <= 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_PORT)
	}
	if len(authApplicationPassword) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_AUTH_APPLICATION_PASSWORD)
	}

	userDirectory := user_directory.New()
	applicationDirectory := application_directory.New(api.ApplicationPassword(authApplicationPassword))
	applicationCheck := application_check.New(applicationDirectory.Check)

	checkHandler := check.New()
	accessDeniedHandler := access_denied.New()
	loginHandler := filter.New(applicationCheck.Check, login.New(userDirectory, applicationDirectory).ServeHTTP, accessDeniedHandler.ServeHTTP)
	applicationCreatorHandler := filter.New(applicationCheck.Check, application_creator.New().ServeHTTP, accessDeniedHandler.ServeHTTP)

	handler := router.New(checkHandler.ServeHTTP, loginHandler.ServeHTTP, applicationCreatorHandler.ServeHTTP)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
