package init

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	logger "github.com/vs-uulm/ztsfc_http_logger"
	"github.com/vs-uulm/ztsfc_http_service/internal/app/config"
)

// InitSysLoggerParams() sets the default values for a system logger
func InitSysLoggerParams() {
	// Set a default logging level.
	// The level "info" is necessary to see the messages
	// of http.Server and httputil.ReverseProxy ErrorLogs.
	if config.Config.SysLogger.LogLevel == "" {
		config.Config.SysLogger.LogLevel = "info"
	}

	// Set a default log messages destination
	if config.Config.SysLogger.LogFilePath == "" {
		config.Config.SysLogger.LogFilePath = "stdout"
	}

	// Set a default log messages JSON formatter
	if config.Config.SysLogger.LogFormatter == "" {
		config.Config.SysLogger.LogFormatter = "json"
	}
}

// InitServiceParams() initializes the 'service' section of the config file
// and loads the Service certificates.
func InitServiceParams(sysLogger *logger.Logger) error {
	var err error
	fields := ""

	if (config.Config.Service == config.ServiceT{}) {
		return fmt.Errorf("init: InitServiceParams(): the section 'service' is empty")
	}

	if config.Config.Service.ListenAddr == "" {
		fields += "listen_addr,"
	}

	if config.Config.Service.CertShownByServiceToClients == "" {
		fields += "cert_shown_by_service_to_clients,"
	}

	if config.Config.Service.PrivkeyForCertShownByServiceToClients == "" {
		fields += "privkey_for_cert_shown_by_service_to_clients,"
	}

	if config.Config.Service.CertServiceAccepts == "" {
		fields += "cert_service_accepts,"
	}

	if config.Config.Service.Mode == "" {
		fields += "direct,"
	}

	if fields != "" {
		return fmt.Errorf("init: InitServiceParams(): in the section 'service' the following required fields are missed: '%s'", strings.TrimSuffix(fields, ","))
	}

	// Preload service X509KeyPair and write it to config
	config.Config.X509KeyPairShownByService, err = loadX509KeyPair(sysLogger,
		config.Config.Service.CertShownByServiceToClients, config.Config.Service.PrivkeyForCertShownByServiceToClients, "service", "")
	if err != nil {
		return err
	}

	// Preload CA certificate and append it to cert pool
	err = loadCACertificate(sysLogger, config.Config.Service.CertServiceAccepts, "service", config.Config.CAcertPoolPepAcceptsFromExt)
	if err != nil {
		return err
	}

	return nil
}

// LoadX509KeyPair() unifies the loading of X509 key pairs for different components
func loadX509KeyPair(sysLogger *logger.Logger, certfile, keyfile, componentName, certAttr string) (tls.Certificate, error) {
	keyPair, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("init: loadX509KeyPair(): loading %s X509KeyPair for %s from '%s' and '%s' - FAIL: %v",
			certAttr, componentName, certfile, keyfile, err)
	}
	sysLogger.Debugf("init: loadX509KeyPair(): loading %s X509KeyPair for %s from '%s' and '%s' - OK", certAttr, componentName, certfile, keyfile)
	return keyPair, nil
}

// function unifies the loading of CA certificates for different components
func loadCACertificate(sysLogger *logger.Logger, certfile string, componentName string, certPool *x509.CertPool) error {
	// Read the certificate file content
	caRoot, err := ioutil.ReadFile(certfile)
	if err != nil {
		return fmt.Errorf("init: loadCACertificate(): loading %s CA certificate from '%s' - FAIL: %w", componentName, certfile, err)
	}
	sysLogger.Debugf("init: loadCACertificate(): loading %s CA certificate from '%s' - OK", componentName, certfile)

	// Return error if provided certificate is nil
	if certPool == nil {
		return errors.New("provided certPool is nil")
	}

	// Append a certificate to the pool
	certPool.AppendCertsFromPEM(caRoot)
	return nil
}

func SetupCloseHandler(logger *logger.Logger) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Debug("- 'Ctrl + C' was pressed in the Terminal. Terminating...")
		logger.Terminate()
		os.Exit(0)
	}()
}
