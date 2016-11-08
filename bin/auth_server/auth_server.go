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
	parameterPort                    = "port"
	parameterAuthApplicationName     = "auth-application-name"
	parameterAuthApplicationPassword = "auth-application-password"
	parameterLedisdbAddress          = "ledisdb-address"
	parameterLedisdbPassword         = "ledisdb-password"
	parameterPrefix                  = "prefix"
)

var (
	portPtr                    = flag.Int(parameterPort, 8080, "port")
	authApplicationNamePtr     = flag.String(parameterAuthApplicationName, "auth", "auth application name")
	authApplicationPasswordPtr = flag.String(parameterAuthApplicationPassword, "", "auth application password")
	ledisdbAddressPtr          = flag.String(parameterLedisdbAddress, "", "ledisdb address")
	ledisdbPasswordPtr         = flag.String(parameterLedisdbPassword, "", "ledisdb password")
	prefixPtr                  = flag.String(parameterPrefix, "", "prefix")
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
		if _, err := factory.ApplicationService().CreateApplicationWithPassword(config.ApplicationName, config.ApplicationPassword); err != nil {
			glog.Warningf("create auth application failed: %v", err)
			return
		}
	}()

	glog.V(2).Infof("start server")
	return gracehttp.Serve(factory.HttpServer())
}

func createConfig() model.Config {
	return model.Config{
		HttpPort:            model.Port(*portPtr),
		HttpPrefix:          model.Prefix(*prefixPtr),
		ApplicationName:     model.ApplicationName(*authApplicationNamePtr),
		ApplicationPassword: model.ApplicationPassword(*authApplicationPasswordPtr),
		LedisdbAddress:      model.LedisdbAddress(*ledisdbAddressPtr),
		LedisdbPassword:     model.LedisdbPassword(*ledisdbPasswordPtr),
	}
}
