package main

import (
	"runtime"

	"github.com/bborbe/auth/factory"
	"github.com/bborbe/auth/model"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/redis_client/ledis"
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
)

var (
	portPtr                    = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "port")
	authApplicationPasswordPtr = flag.String(PARAMETER_AUTH_APPLICATION_PASSWORD, "", "auth application password")
	ledisdbAddressPtr          = flag.String(PARAMETER_LEDISDB_ADDRESS, "", "ledisdb address")
	ledisdbPasswordPtr         = flag.String(PARAMETER_LEDISDB_PASSWORD, "", "ledisdb password")
	prefixPtr                  = flag.String(PARAMETER_PREFIX, "", "prefix")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {

	config := createConfig()
	if err := config.Validate(); err != nil {
		return err
	}

	factory := factory.New(config, ledis.New(config.LedisdbAddress.String(), config.LedisdbPassword.String()))

	go func() {
		if _, err := factory.ApplicationService().CreateApplicationWithPassword(model.AUTH_APPLICATION_NAME, config.ApplicationPassword); err != nil {
			glog.Warningf("create auth application failed: %v", err)
			return
		}
	}()

	glog.V(2).Infof("start server")
	return gracehttp.Serve(factory.HttpServer())
}

func createConfig() model.Config {
	return model.Config{
		Port:                model.Port(*portPtr),
		Prefix:              model.Prefix(*prefixPtr),
		ApplicationPassword: model.ApplicationPassword(*authApplicationPasswordPtr),
		LedisdbAddress:      model.LedisdbAddress(*ledisdbAddressPtr),
		LedisdbPassword:     model.LedisdbPassword(*ledisdbPasswordPtr),
	}
}
