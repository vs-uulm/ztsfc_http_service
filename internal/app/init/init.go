package init

import (
	"fmt"
    "crypto/x509"
    "crypto/tls"
	"strings"
    "net"

    "gopkg.in/ldap.v2"

    gct "github.com/leobrada/golang_convenience_tools"
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

func InitConfig(sysLogger *logger.Logger) error {
    if err := initService(); err != nil {
        return fmt.Errorf("init: InitConfig(): %v", err)
    }

    if config.Config.Service.Mode == "direct" {
        if err := initBasicAuth(sysLogger); err != nil {
            return fmt.Errorf("init: InitConfig(): %v", err)
        }
    } else {
        sysLogger.Infof("init: InitConfig(): 'service' part is skipped since service is not running in mode 'direct'")
    }

    return nil
}

// InitServiceParams() initializes the 'service' section of the config file
// and loads the Service certificates.
func initService() error {
	var err error
	fields := ""

	//if (config.Config.Service == config.ServiceT{}) {
	//	return fmt.Errorf("init: InitServiceParams(): the section 'service' is empty")
	//}

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
    config.Config.Service.CAcertPoolServiceAcceptsFromExt = x509.NewCertPool()
    if err = gct.LoadCACertificate(config.Config.Service.CertServiceAccepts, config.Config.Service.CAcertPoolServiceAcceptsFromExt); err != nil {
        return fmt.Errorf("initService():  error loading certificate service accepts from clients: %w", err)
    }

    config.Config.Service.X509KeyPairShownByService, err = gct.LoadX509KeyPair(config.Config.Service.CertShownByServiceToClients,
        config.Config.Service.PrivkeyForCertShownByServiceToClients)

	return nil
}

func initBasicAuth(sysLogger *logger.Logger) error {
    err := initSession(sysLogger)
    if err != nil {
        return err
    }

    err = initLdap(sysLogger)
    if err != nil {
        return err
    }

    err = initPerimeter(sysLogger)
    return err
}

func initSession(sysLogger *logger.Logger) error {
        var err error
        fields := ""

        if config.Config.BasicAuth.Session.Path_to_jwt_pub_key == "" {
                fields += "path_to_jwt_pub_key,"
        }
        sysLogger.Debugf("init: initSession(): JWT Public Key path: '%s'", config.Config.BasicAuth.Session.Path_to_jwt_pub_key)

        if config.Config.BasicAuth.Session.Path_to_jwt_signing_key == "" {
                fields += "path_to_jwt_signing_key,"
        }
        sysLogger.Debugf("init: initSession(): JWT Signing Key path: '%s'", config.Config.BasicAuth.Session.Path_to_jwt_signing_key)

        if fields != "" {
                return fmt.Errorf("init: initSession(): in the section 'session' the following required fields are missed: '%s'", strings.TrimSuffix(fields, ","))
        }

        config.Config.BasicAuth.Session.JwtPubKey, err = gct.ParseRsaPublicKeyFromPemFile(config.Config.BasicAuth.Session.Path_to_jwt_pub_key)
        if err != nil {
                return err
        }

        config.Config.BasicAuth.Session.MySigningKey, err = gct.ParseRsaPrivateKeyFromPemFile(config.Config.BasicAuth.Session.Path_to_jwt_signing_key)
        if err != nil {
                return err
        }

        return nil
}

// InitLdapParams() initializes the 'ldap' section of the config file.
func initLdap(sysLogger *logger.Logger) error {
	var err error
	fields := ""

	// TODO: Check if the field make sense as well!
	if config.Config.BasicAuth.Ldap.Base == "" {
		fields += "base,"
	}

	// TODO: Check if the field make sense as well!
	if config.Config.BasicAuth.Ldap.Host == "" {
		fields += "host,"
	}

	// TODO: Check if the field make sense as well!
	if config.Config.BasicAuth.Ldap.Port <= 0 {
		fields += "port,"
	}

	// TODO: Check if the field make sense as well!
	//if config.Config.BasicAuth.Ldap.BindDN == "" {
	//	fields += "bind_dn,"
	//}

	// TODO: Check if the field make sense as well!
	//if config.Config.BasicAuth.Ldap.BindPassword == "" {
	//	fields += "bind_password,"
	//}

	// TODO: Check if the field make sense as well!
	if config.Config.BasicAuth.Ldap.UserFilter == "" {
		fields += "user_filter,"
	}

	// TODO: Check if the field make sense as well!
	if config.Config.BasicAuth.Ldap.ReadonlyDN == "" {
		fields += "readonly_dn,"
	}

	// TODO: Check if the field make sense as well!
	if config.Config.BasicAuth.Ldap.ReadonlyPW == "" {
		fields += "readonly_pw,"
	}

	// TODO: Check if the field make sense as well!
	//if config.Config.BasicAuth.Ldap.GroupFilter == "" {
	//	fields += "group_filter,"
	//}

	// TODO: Check if the field make sense as well!
	if len(config.Config.BasicAuth.Ldap.Attributes) == 0 {
		fields += "attributes,"
	}

	if fields != "" {
		return fmt.Errorf("init: InitLdap(): in the section 'ldap' the following required fields are missed: '%s'", strings.TrimSuffix(fields, ","))
	}

	// Preload X509KeyPair and write it to config
	config.Config.BasicAuth.Ldap.X509KeyPairShownByServiceToLdap, err = gct.LoadX509KeyPair(config.Config.BasicAuth.Ldap.CertShownByServiceToLdap, config.Config.BasicAuth.Ldap.PrivkeyForCertShownByServiceToLdap)
	if err != nil {
		return err
	}

	// Preload CA certificate and append it to cert pool
    if err = gct.LoadCACertificate(config.Config.BasicAuth.Ldap.CertServiceAcceptsShownByLdap, config.Config.Service.CAcertPoolServiceAcceptsFromExt); err != nil {
        return fmt.Errorf("initLdap():  error loading certificate service accepts from clients: %w", err)
    }

	// Create an LDAP connection
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{config.Config.BasicAuth.Ldap.X509KeyPairShownByServiceToLdap},
		RootCAs:      config.Config.Service.CAcertPoolServiceAcceptsFromExt,
		ServerName:   config.Config.BasicAuth.Ldap.Host,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
        InsecureSkipVerify: true,
	}

	config.Config.BasicAuth.Ldap.LdapConn, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", config.Config.BasicAuth.Ldap.Host,
        config.Config.BasicAuth.Ldap.Port), tlsConf)
	if err != nil {
		return fmt.Errorf("init: initLdap(): unable to connect to the LDAP server: %s", err.Error())
	}

	return nil
}

func initPerimeter(sysLogger *logger.Logger) error {
    // Iterates over all trusted locations (for each resource) and tries to extract the IPNet from it
    for _, location := range config.Config.BasicAuth.Perimeter.TrustedLocations {
        _, ipnet, err := net.ParseCIDR(location)
        if err != nil {
            return fmt.Errorf("init: InitResourcesParams(): %s is not in valid CIDR network notation: %v",
                location, err)
        }
        config.Config.BasicAuth.Perimeter.TrustedIPNetworks = append(config.Config.BasicAuth.Perimeter.TrustedIPNetworks, ipnet)
    }

    return nil
}

//func SetupCloseHandler(logger *logger.Logger) {
//	c := make(chan os.Signal)
//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//	go func() {
//		<-c
//		logger.Debug("- 'Ctrl + C' was pressed in the Terminal. Terminating...")
//		logger.Terminate()
//		os.Exit(0)
//	}()
//}
