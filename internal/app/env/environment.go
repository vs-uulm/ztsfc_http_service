package env

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"

	logwriter "github.com/vs-uulm/ztsfc_http_service/internal/app/logwriter"
//    github.com/vs-uulm/ztsfc_http_pdp/internal/app/
)

var Config Config_t

type Config_t struct {
	Sni                                         string   `yaml:"sni"`
	Listen_addr                                 string   `yaml:"listen_addr"`
	Certs_service_accepts_when_presented_by_ext []string `yaml:"certs_service_accepts_when_presented_by_ext"`
	Certs_service_accepts_when_presented_by_int []string `yaml:"certs_service_accepts_when_presented_by_int"`

	Cert_presented_by_service_to_ext_matching_sni string `yaml:"cert_presented_by_service_to_ext_matching_sni"`
	Privkey_for_cert_presented_by_service_to_ext  string `yaml:"privkey_for_cert_presented_by_service_to_ext"`
	Cert_presented_by_service_to_int              string `yaml:"cert_presented_by_service_to_int"`
	Privkey_for_cert_presented_by_service_to_int  string `yaml:"privkey_for_cert_presented_by_service_to_int"`

	CA_cert_pool_service_accepts_when_presented_by_ext *x509.CertPool
	CA_cert_pool_service_accepts_when_presented_by_int *x509.CertPool
	X509KeyPair_presented_by_service_to_ext            tls.Certificate
	X509KeyPair_presented_by_service_to_int            tls.Certificate
}

// Parses a configuration yaml file into the global Config variable
func LoadConfig(configPath string, lw *logwriter.LogWriter) (err error) {
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		lw.Logger.WithFields(logrus.Fields{"type": "system"}).Fatalf("Open configuration file error: %v", err)
	} else {
		lw.Logger.WithFields(logrus.Fields{"type": "system"}).Debugf("Configuration file %s exists and is readable", configPath)
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	err = d.Decode(&Config)
	if err != nil {
		lw.Logger.WithFields(logrus.Fields{"type": "system"}).Fatalf("Configuration yaml-->go decoding error: %v", err)
	} else {
		lw.Logger.WithFields(logrus.Fields{"type": "system"}).Debugf("Configuration has been successfully decoded")
	}

	return
}
