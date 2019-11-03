package main

import (
	"context"
	"net/http"

	"github.com/qaisjp/gosign"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// API contains all the dependencies of the API server
type API struct {
	Config *Config
	Log    *logrus.Logger
	GoSign *gosign.Client
	Gin    *gin.Engine

	Server *http.Server
	Tokens map[string]string
}

// Start binds the API and starts listening.
func (a *API) Start() error {
	a.Server = &http.Server{
		Addr:    a.Config.Address,
		Handler: a.Gin,
	}
	return a.Server.ListenAndServe()
}

// Shutdown shuts down the API
func (a *API) Shutdown(ctx context.Context) error {
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// NewAPI sets up a new API module.
func NewAPI(
	conf *Config,
	log *logrus.Logger,
	gosign *gosign.Client,
	tokens map[string]string,
) *API {

	router := gin.Default()

	a := &API{
		Config: conf,
		Log:    log,
		GoSign: gosign,
		Gin:    router,
		Tokens: tokens,
	}

	router.GET("/cosign/valid", a.Valid)
	router.GET("/check/:token_name/:token_key", a.Check)

	return a
}
