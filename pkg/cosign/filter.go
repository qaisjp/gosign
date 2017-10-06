package cosign

import (
	"crypto/tls"

	"github.com/compsoc-edinburgh/bi-dice-api/pkg/config"
	cosign "github.com/qaisjp/go-cosign"
)

func NewFilter(cfg config.CoSignConfig) (*cosign.Filter, error) {
	filter, err := cosign.Dial(&cosign.Config{
		Address: cfg.DaemonAddress,
		Service: cfg.Service,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: cfg.Insecure,
		},
	})
	if err != nil {
		return nil, err
	}

	return filter, nil
}
