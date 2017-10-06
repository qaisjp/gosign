package cosign

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/compsoc-edinburgh/bi-dice-api/pkg/config"
	"github.com/pkg/errors"
	cosign "github.com/qaisjp/go-cosign"
)

func NewFilter(cfg config.CoSignConfig) (*cosign.Filter, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not read certfile+keyfile")
	}

	pool := x509.NewCertPool()
	// fmt.Println("Append ok? ", pool.AppendCertsFromPEM([]byte(pemCerts)))

	filter, err := cosign.Dial(&cosign.Config{
		Address: cfg.DaemonAddress,
		Service: cfg.Service,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: cfg.Insecure,
			ServerName:         cfg.ServerName,
			Certificates:       []tls.Certificate{cert},
			ClientCAs:          pool,
		},
	})
	if err != nil {
		return nil, err
	}

	return filter, nil
}

// var pemCerts = `content snipped`
