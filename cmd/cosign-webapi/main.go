package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/koding/multiconfig"
)

func main() {
	var err error

	m := multiconfig.NewWithPath(os.Getenv("config"))
	cfg := &Config{}
	m.MustLoad(cfg)

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	logger := logrus.StandardLogger()
	logger.Level = logLevel

	logger.WithFields(logrus.Fields{
		"module": "init",
	}).Info("Starting up dice-api")

	// Make sure no token has the same name
	tokens := map[string]string{}
	for _, token := range cfg.Tokens {
		// check tok exists
		_, ok := tokens[token.Name]
		if ok {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"name":   token.Name,
			}).Info("multiple tokens exist with the same name")

			return
		}

		// ok, doesn't exist, append
		tokens[token.Name] = token.Key
	}

	// Initialize the cosign client
	filter, err := newGosignClient(cfg.CoSign)

	addr := net.JoinHostPort(cfg.CoSign.DaemonHost, cfg.CoSign.DaemonPort)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"error":  err.Error(),
			"addr":   addr,
		}).Fatal("Unable to connect to the CoSign daemon")
		return
	}

	logger.WithFields(logrus.Fields{
		"module": "init",
		"addr":   addr,
	}).Info("Connected to the CoSign daemon")

	api := NewAPI(
		cfg,
		logger,
		filter,
		tokens,
	)

	go func() {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"bind":   cfg.Address,
		}).Info("Starting the API server")

		if err := api.Start(); err != nil {
			logger.WithFields(logrus.Fields{
				"module": "init",
				"error":  err.Error(),
			}).Fatal("API server failed")
		}
	}()

	// Create a new signal receiver
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Watch for a signal
	<-sc

	// ugly thing to stop ^C from killing alignment
	logger.Out.Write([]byte("\r\n"))

	if err := filter.Quit(); err != nil {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"error":  err.Error(),
		}).Fatal("Failed to close the API server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := api.Shutdown(ctx); err != nil {
		logger.WithFields(logrus.Fields{
			"module": "init",
			"error":  err.Error(),
		}).Fatal("Failed to close the API server")
	}

	logger.WithFields(logrus.Fields{
		"module": "init",
	}).Info("dice-api has shut down.")
}
