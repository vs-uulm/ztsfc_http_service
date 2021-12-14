package main

import (
	"crypto/x509"
	"flag"
	"net/http"
	env "local.com/leobrada/ztsfc_http_pep/env"
	sinit "local.com/leobrada/ztsfc_http_pep/init"
	router "local.com/leobrada/ztsfc_http_pep/router"
	logwriter "local.com/leobrada/ztsfc_http_pep/logwriter"
	"github.com/sirupsen/logrus"
)

var (
    mode string
    file bool
	conf_file_path string
	log_file_path string
	log_level string
	ifTextFormatter bool

	// An instance of logwriter based on logrus
	lw *logwriter.LogWriter
)

func init() {
    flag.StringVar(&mode, "m", "direct", "Specifies the operation mode; either 'direct' or 'pep'")
    flag.BoolVar(&file, "f", false, "Specifies if the service should serve a file or not")
	flag.StringVar(&conf_file_path, "c", "./conf.yml", "Path to user defined yml config file")
	flag.StringVar(&log_file_path, "l", "./service.log", "Path to log file")
    flag.StringVar(&log_level, "log-level", "error", "Log level from the next set: debug, info, warning, error")
    flag.BoolVar(&ifTextFormatter, "text", false, "Use a text format instead of JSON to log messages")

	// Operating input parameters
	flag.Parse()

	lw = logwriter.New(log_file_path, log_level, ifTextFormatter)
	sysLogger := lw.Logger.WithFields(logrus.Fields{"type": "system"})
	sinit.SetupCloseHandler(lw)

	// Loading all config parameter from config file defined in "conf_file_path"
	err := env.LoadConfig(conf_file_path, lw)
	if err != nil {
		sysLogger.Fatalf("Loading logger configuration from %s - ERROR: %v", conf_file_path, err)
	} else {
		sysLogger.Debugf("Loading logger configuration from %s - OK", conf_file_path)
	}

	// Create Certificate Pools for the CA certificates used by the PEP
	env.Config.CA_cert_pool_service_accepts_when_presented_by_ext = x509.NewCertPool()
	env.Config.CA_cert_pool_service_accepts_when_presented_by_int = x509.NewCertPool()

	// Load all CA certificates
	err = sinit.LoadServiceCerts(lw)
	if err != nil {
		sysLogger.Fatalf("Loading service certificates - ERROR: %v", err)
	} else {
		sysLogger.WithFields(logrus.Fields{"type":"system"}).Debug("Loading service certificates - OK")
	}
}

func main() {
	// Create new Service router
	r, err := router.NewRouter(lw, mode, file)
	if err != nil {
		lw.Logger.Fatalf("Fatal error during new router creation: %v", err)
	} else {
		lw.Logger.WithFields(logrus.Fields{"type":"system"}).Debug("New router is successfully created")
	}

	http.Handle("/", r)

	err = r.ListenAndServeTLS()
	if err != nil {
		lw.Logger.Fatalf("ListenAndServeTLS Fatal Error: %v", err)
	}
}
