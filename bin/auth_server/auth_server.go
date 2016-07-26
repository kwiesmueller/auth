package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/bborbe/auth/handler_creator"
	flag "github.com/bborbe/flagenv"
	io_util "github.com/bborbe/io/util"
	"github.com/bborbe/log"
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
	PARAMETER_PREFIX                    = "prefix"
	PARAMETER_ROOT                      = "root"
)

var (
	logLevelPtr                = flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	portPtr                    = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
	ledisdbAddressPtr          = flag.String(PARAMETER_LEDISDB_ADDRESS, "", "ledisdb address")
	ledisdbPasswordPtr         = flag.String(PARAMETER_LEDISDB_PASSWORD, "", "ledisdb password")
	prefixPtr                  = flag.String(PARAMETER_PREFIX, "", "prefix")
	rootPtr                    = flag.String(PARAMETER_ROOT, "", "root")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	server, err := createServer(
		*portPtr,
		*prefixPtr,
		*rootPtr,
		*authApplicationPasswordPtr,
		*ledisdbAddressPtr,
		*ledisdbPasswordPtr,
	)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debugf("start server")
	gracehttp.Serve(server)
}

func createServer(
	port int,
	prefix string,
	documentRoot string,
	authApplicationPassword string,
	ledisdbAddress string,
	ledisdbPassword string,
) (*http.Server, error) {
	logger.Debugf("create server with port: %d", port)
	if port <= 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_PORT)
	}
	if len(documentRoot) == 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_ROOT)
	}
	documentRoot, err := io_util.NormalizePath(documentRoot)
	if err != nil {
		return nil, err
	}
	if len(ledisdbAddress) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_LEDISDB_ADDRESS)
	}
	handlerCreator := handler_creator.New()
	handler, err := handlerCreator.CreateHandler(
		prefix,
		documentRoot,
		authApplicationPassword,
		ledisdbAddress,
		ledisdbPassword,
	)
	if err != nil {
		return nil, err
	}
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
