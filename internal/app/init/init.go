package init

import (
	"github.com/vs-uulm/ztsfc_http_service/internal/app/config"
	logwriter "github.com/vs-uulm/ztsfc_http_service/internal/app/logwriter"
	"crypto/tls"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"github.com/sirupsen/logrus"
)

func LoadServiceCerts(lw *logwriter.LogWriter) error {
    var ca []byte
    var err error

    config.Config.X509KeyPair_presented_by_service_to_ext, err = tls.LoadX509KeyPair(
        config.Config.Cert_presented_by_service_to_ext_matching_sni, config.Config.Privkey_for_cert_presented_by_service_to_ext)
    if err != nil {
        lw.Logger.WithFields(logrus.Fields{"type":"system"}).Fatalf("Critical error when loading X509KeyPair_presented_by_service_to_ext")
    } else {
        lw.Logger.WithFields(logrus.Fields{"type":"system"}).Debugf("X509KeyPair_presented_by_service_to_ext successfully loaded")
    }

    config.Config.X509KeyPair_presented_by_service_to_int, err = tls.LoadX509KeyPair(
        config.Config.Cert_presented_by_service_to_int, config.Config.Privkey_for_cert_presented_by_service_to_int)
    if err != nil {
        lw.Logger.WithFields(logrus.Fields{"type":"system"}).Fatalf("Critical error when loading X509KeyPair_presented_by_service_to_int")
    } else {
        lw.Logger.WithFields(logrus.Fields{"type":"system"}).Debugf("X509KeyPair_presented_by_service_to_int successfully loaded")
    }

	// Read CA certs used for signing ext certs (especially client certs) and are accepted by the service
	for _, acceptedExtCert := range config.Config.Certs_service_accepts_when_presented_by_ext {
		ca, err = ioutil.ReadFile(acceptedExtCert)
		if err != nil {
			lw.Logger.WithFields(logrus.Fields{"type":"system"}).Fatalf("Loading external CA certificate from %s error", acceptedExtCert)
		} else {
			lw.Logger.WithFields(logrus.Fields{"type":"system"}).Debugf("External CA certificate from %s is successfully loaded", acceptedExtCert)
		}
		// Append a certificate to the pool
		config.Config.CA_cert_pool_service_accepts_when_presented_by_ext.AppendCertsFromPEM(ca)
	}

	// Read CA certs used for signing int certs and are accepted by the service
	for _, acceptedIntCert := range config.Config.Certs_service_accepts_when_presented_by_int {
		ca, err = ioutil.ReadFile(acceptedIntCert)
		if err != nil {
			lw.Logger.WithFields(logrus.Fields{"type":"system"}).Fatalf("Loading external CA certificate from %s error", acceptedIntCert)
		} else {
			lw.Logger.WithFields(logrus.Fields{"type":"system"}).Debugf("External CA certificate from %s is successfully loaded", acceptedIntCert)
		}
		// Append a certificate to the pool
		config.Config.CA_cert_pool_service_accepts_when_presented_by_int.AppendCertsFromPEM(ca)
	}

    return err
}

func SetupCloseHandler(lw *logwriter.LogWriter) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		lw.Logger.WithFields(logrus.Fields{"type":"system"}).Debug("- Ctrl+C pressed in Terminal. Terminating...")
		lw.Terminate()
		os.Exit(0)
	}()
}
