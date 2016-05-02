package main

import (
	"fmt"
	"net/http"
	"os"

	auth_check "github.com/bborbe/auth/check"
	auth_login "github.com/bborbe/auth/login"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/log"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
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
	logger.Debugf("create server with port: %d", port)

	check := auth_check.New()
	login := auth_login.New()

	router := mux.NewRouter()
	router.Path("/healthz").Methods("GET").HandlerFunc(check.ServeHTTP)
	router.Path("/readiness").Methods("GET").HandlerFunc(check.ServeHTTP)
	router.Path("/login").Methods("POST").HandlerFunc(login.ServeHTTP)

	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}, nil
}
