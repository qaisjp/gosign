package api

import (
	"github.com/compsoc-edinburgh/bi-dice-api/pkg/api/backend"
	"github.com/compsoc-edinburgh/bi-dice-api/pkg/api/base"
	"github.com/compsoc-edinburgh/bi-dice-api/pkg/api/frontend"
	"github.com/compsoc-edinburgh/bi-dice-api/pkg/config"
	"github.com/gin-gonic/gin"
	cosign "github.com/qaisjp/go-cosign"
	"github.com/sirupsen/logrus"
)

// NewAPI sets up a new API module.
func NewAPI(
	conf *config.Config,
	log *logrus.Logger,
	filter *cosign.Filter,
) *base.API {

	router := gin.Default()

	a := &base.API{
		Config: conf,
		Log:    log,
		Filter: filter,
		Gin:    router,
	}

	frontend := frontend.Impl{API: a}
	router.GET("/cosign/valid", frontend.Valid)

	backend := backend.Impl{API: a}
	router.GET("/check/:username/:value/:login_cookie", backend.Check)

	return a
}
