// Package config reads the config file and parses it to go data structures.
package config

import (
	"crypto/tls"
    "crypto/rsa"
	"crypto/x509"
    "net"

    "gopkg.in/ldap.v2"
)

// Config contains all input from the config file and is is globally accessible
var Config ConfigT

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

	CAcertPoolServiceAcceptsFromExt *x509.CertPool
	X509KeyPairShownByService   tls.Certificate
}

type BasicAuthT struct {
    Session SessionT `yaml:"session"`
    Ldap       LdapT       `yaml:"ldap"`
    Perimeter PerimeterT `yaml:"perimeter"`
}

type SessionT struct {
    Path_to_jwt_pub_key     string `yaml:"path_to_jwt_pub_key"`
    Path_to_jwt_signing_key string `yaml:"path_to_jwt_signing_key"`
    JwtPubKey               *rsa.PublicKey
    MySigningKey            *rsa.PrivateKey
}

// The struct LdapT is for parsing the section 'ldap' of the config file.
type LdapT struct {
        Base         string   `yaml:"base"`
        Host         string   `yaml:"host"`
        Port         int      `yaml:"port"`
        UseSSL       bool     `yaml:"use_ssl"`
        BindDN       string   `yaml:"bind_dn"`
        BindPassword string   `yaml:"bind_password"`
        ReadonlyDN   string   `yaml:"readonly_dn"`
        ReadonlyPW   string   `yaml:"readonly_pw"`
        UserFilter   string   `yaml:"user_filter"`
        GroupFilter  string   `yaml:"group_filter"`
        Attributes   []string `yaml:"attributes"`

        CertShownByServiceToLdap           string `yaml:"cert_shown_by_service_to_ldap"`
        PrivkeyForCertShownByServiceToLdap string `yaml:"privkey_for_cert_shown_by_service_to_ldap"`
        CertServiceAcceptsShownByLdap      string `yaml:"cert_service_accepts_shown_by_ldap"`
        X509KeyPairShownByServiceToLdap    tls.Certificate
        LdapConn                       *ldap.Conn
}

type PerimeterT struct {
    ApplyPerimeter bool `yaml:"apply_perimeter"`
    TrustedLocations   []string `yaml:"trusted_locations"`
    TrustedIPNetworks []*net.IPNet
}

// ConfigT struct is for parsing the basic structure of the config file
type ConfigT struct {
	SysLogger                   sysLoggerT `yaml:"system_logger"`
	Service                     ServiceT   `yaml:"service"`
    BasicAuth  BasicAuthT  `yaml:"basic_auth"`
}
