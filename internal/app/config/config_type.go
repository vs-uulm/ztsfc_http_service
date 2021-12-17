// Package config reads the config file and parses it to go data structures.
package config

import (
	"crypto/tls"
	"crypto/x509"
)

type sysLoggerT struct {
	LogLevel     string `yaml:"system_logger_logging_level"`
	LogFilePath  string `yaml:"system_logger_destination"`
	LogFormatter string `yaml:"system_logger_format"`
}

// The struct ServiceT is for parsing the section 'service' of the config file.
type ServiceT struct {
	ListenAddr                            string `yaml:"listen_addr"`
	CertShownByServiceToClients           string `yaml:"cert_shown_by_service_to_clients"`
	PrivkeyForCertShownByServiceToClients string `yaml:"privkey_for_cert_shown_by_service_to_clients"`
	CertServiceAccepts                    string `yaml:"cert_service_accepts"`
	Mode                                  string `yaml:"mode"`
	File                                  bool   `yaml:"file"`
}

// ConfigT struct is for parsing the basic structure of the config file
type ConfigT struct {
	SysLogger                   sysLoggerT `yaml:"system_logger"`
	Service                     ServiceT   `yaml:"service"`
	CAcertPoolPepAcceptsFromExt *x509.CertPool
	X509KeyPairShownByService   tls.Certificate
}

// Config contains all input from the config file and is is globally accessible
var Config ConfigT
