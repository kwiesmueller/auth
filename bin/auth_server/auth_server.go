package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/bborbe/auth/handler_creator"
	flag "github.com/bborbe/flagenv"
	debug_handler "github.com/bborbe/http_handler/debug"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
)

const (
	DEFAULT_PORT                        = 8080
	PARAMETER_PORT                      = "port"
	PARAMETER_AUTH_APPLICATION_PASSWORD = "auth-application-password"
	PARAMETER_LEDISDB_ADDRESS           = "ledisdb-address"
	PARAMETER_LEDISDB_PASSWORD          = "ledisdb-password"
	PARAMETER_PREFIX                    = "prefix"
	PARAMETER_DEBUG                     = "debug"
)

var (
	portPtr                    = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
	ledisdbAddressPtr          = flag.String(PARAMETER_LEDISDB_ADDRESS, "", "ledisdb address")
	ledisdbPasswordPtr         = flag.String(PARAMETER_LEDISDB_PASSWORD, "", "ledisdb password")
	prefixPtr                  = flag.String(PARAMETER_PREFIX, "", "prefix")
	debugPtr                   = flag.Bool(PARAMETER_DEBUG, false, "debug")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := do(
		*portPtr,
		*debugPtr,
		*prefixPtr,
		*authApplicationPasswordPtr,
		*ledisdbAddressPtr,
		*ledisdbPasswordPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	port int,
	debug bool,
	prefix string,
	authApplicationPassword string,
	ledisdbAddress string,
	ledisdbPassword string,
) error {
	server, err := createServer(
		port,
		debug,
		prefix,
		authApplicationPassword,
		ledisdbAddress,
		ledisdbPassword,
	)
	if err != nil {
		return err
	}
	glog.V(2).Infof("start server")
	return gracehttp.Serve(server)
}

func createServer(
	port int,
	debug bool,
	prefix string,
	authApplicationPassword string,
	ledisdbAddress string,
	ledisdbPassword string,
) (*http.Server, error) {
	glog.V(2).Infof("create server with port: %d", port)
	if port <= 0 {
		return nil, fmt.Errorf("parameter %s invalid", PARAMETER_PORT)
	}
	if len(ledisdbAddress) == 0 {
		return nil, fmt.Errorf("parameter %s missing", PARAMETER_LEDISDB_ADDRESS)
	}
	handlerCreator := handler_creator.New()
	handler, err := handlerCreator.CreateHandler(
		prefix,
		authApplicationPassword,
		ledisdbAddress,
		ledisdbPassword,
	)
	if err != nil {
		return nil, err
	}

	if debug {
		handler = debug_handler.New(handler)
	}

	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
