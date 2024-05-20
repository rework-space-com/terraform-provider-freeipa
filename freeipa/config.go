package freeipa

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
)

// Config is the configuration parameters for the FreeIPA API
type Config struct {
	Host               string
	Username           string
	Password           string
	InsecureSkipVerify bool
	CaCertificate      string
}

// Client creates a FreeIPA client scoped to the global API
func (c *Config) Client() (*ipa.Client, error) {
	caCertPool := x509.NewCertPool()

	if c.CaCertificate != "" {
		caCert, err := os.ReadFile(c.CaCertificate)
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.InsecureSkipVerify,
			RootCAs:            caCertPool,
		},
	}

	client, err := ipa.Connect(c.Host, tspt, c.Username, c.Password)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] FreeIPA Client configured for host: %s", c.Host)

	return client, nil
}
