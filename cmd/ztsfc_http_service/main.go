package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/vs-uulm/ztsfc_http_service/internal/app/config"
	logger "github.com/vs-uulm/ztsfc_http_logger"
	confInit "github.com/vs-uulm/ztsfc_http_service/internal/app/init"
	router "github.com/vs-uulm/ztsfc_http_service/internal/app/router"
    yaml "github.com/leobrada/yaml_tools"
)

var (
	sysLogger    *logger.Logger
)

func init() {
	var confFilePath string

	// Operating input parameters
	flag.StringVar(&confFilePath, "c", "./config/conf.yml", "Path to user defined YML config file")
	flag.Parse()

	// Loading all config parameter from config file defined in "confFilePath"
    err := yaml.LoadYamlFile(confFilePath, &config.Config)
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

    if err = confInit.InitConfig(sysLogger); err != nil {
        sysLogger.Fatalf("main: init(): could not initialize Service params: %v", err)
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
