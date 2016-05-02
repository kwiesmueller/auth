package main

import (
	"flag"

	"fmt"
	"net/http"
	"os"

	auth_check "github.com/bborbe/auth/check"
	auth_login "github.com/bborbe/auth/login"
	"github.com/bborbe/log"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
)

var logger = log.DefaultLogger

const (
	DEFAULT_PORT       int = 8080
	PARAMETER_LOGLEVEL     = "loglevel"
	PARAMETER_PORT         = "port"
	PARAMETER_ROOT         = "root"
)

var (
	logLevelPtr = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	portPtr     = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
	rootPtr     = flag.String(PARAMETER_ROOT, "/data", "auth root directory")
)

func main() {
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	server, err := createServer(*portPtr, *rootPtr)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	gracehttp.Serve(server)
}

func createServer(port int, root string) (*http.Server, error) {
	logger.Debugf("create server with port: %d root: %s", port, root)

	check := auth_check.New()
	login := auth_login.New()

	router := mux.NewRouter()
	router.Path("/healthz").Methods("GET").HandlerFunc(check.ServeHTTP)
	router.Path("/readiness").Methods("GET").HandlerFunc(check.ServeHTTP)
	router.Path("/login").Methods("GET").HandlerFunc(login.ServeHTTP)

	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}, nil
}
