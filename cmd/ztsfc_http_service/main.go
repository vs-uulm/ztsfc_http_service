package main

import (
	"crypto/x509"
	"flag"
	"log"
	"net/http"

	logger "github.com/vs-uulm/ztsfc_http_logger"
	"github.com/vs-uulm/ztsfc_http_service/internal/app/config"
	confInit "github.com/vs-uulm/ztsfc_http_service/internal/app/init"
	router "github.com/vs-uulm/ztsfc_http_service/internal/app/router"
)

var (
	confFilePath string
	sysLogger    *logger.Logger
)

func init() {
	var err error

	// Operating input parameters
	flag.StringVar(&confFilePath, "c", "./config/conf.yml", "Path to user defined YML config file")
	flag.Parse()

	// Loading all config parameter from config file defined in "confFilePath"
	err = config.LoadConfig(confFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create an instance of the system logger
	confInit.InitSysLoggerParams()
	sysLogger, err = logger.New(config.Config.SysLogger.LogFilePath,
		config.Config.SysLogger.LogLevel,
		config.Config.SysLogger.LogFormatter,
		logger.Fields{"type": "system"},
	)
	if err != nil {
		log.Fatal(err)
	}
	sysLogger.Debugf("loading logger configuration from '%s' - OK", confFilePath)

	// Create Certificate Pools for the CA certificates used by the service
	config.Config.CAcertPoolPepAcceptsFromExt = x509.NewCertPool()

	// service
	err = confInit.InitServiceParams(sysLogger)
	if err != nil {
		sysLogger.Fatal(err)
	}
}

func main() {
	// Create a new Service router
	serviceRouter, err := router.NewRouter(sysLogger, config.Config.Service.Mode, config.Config.Service.File)
	if err != nil {
		sysLogger.Error(err)
		return
	}
	sysLogger.Debug("main: new router was successfully created")

	http.Handle("/", serviceRouter)

	err = serviceRouter.ListenAndServeTLS()
	if err != nil {
		sysLogger.Error(err)
	}
}
