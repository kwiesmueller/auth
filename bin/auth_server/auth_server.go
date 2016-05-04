package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bborbe/auth/application_directory"
	auth_check "github.com/bborbe/auth/check"
	auth_login "github.com/bborbe/auth/login"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/user_directory"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/log"
	"github.com/facebookgo/grace/gracehttp"
)

var logger = log.DefaultLogger

const (
	DEFAULT_PORT       = 8080
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_PORT     = "port"
)

var (
	logLevelPtr = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	portPtr     = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
)

func main() {
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	server, err := createServer(*portPtr)

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	gracehttp.Serve(server)
}

func createServer(port int) (*http.Server, error) {
	if port <= 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_PORT)
	}
	logger.Debugf("create server with port: %d", port)

	userDirectory := user_directory.New()
	applicationDirectory := application_directory.New()
	check := auth_check.New()
	login := auth_login.New(userDirectory, applicationDirectory)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router.New(check.ServeHTTP, login.ServeHTTP)}, nil
}
