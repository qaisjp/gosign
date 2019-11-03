package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qaisjp/gosign/internal/api/backend"
	"github.com/qaisjp/gosign/internal/api/base"
	"github.com/qaisjp/gosign/internal/api/frontend"
	"github.com/qaisjp/gosign/internal/config"
	"github.com/qaisjp/gosign"
	"github.com/sirupsen/logrus"
)

// NewAPI sets up a new API module.
func NewAPI(
	conf *config.Config,
	log *logrus.Logger,
	gosign *gosign.Client,
	tokens map[string]string,
) *base.API {

	router := gin.Default()

	a := &base.API{
		Config: conf,
		Log:    log,
		GoSign: gosign,
		Gin:    router,
		Tokens: tokens,
	}

	frontend := frontend.Impl{API: a}
	router.GET("/cosign/valid", frontend.Valid)

	backend := backend.Impl{API: a}
	router.GET("/check/:token_name/:token_key", backend.Check)

	return a
}
