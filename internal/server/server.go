package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mchmarny/vul/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	portEnvVar        = "PORT"
	portDefaultVal    = "8080"
	configEnvVar      = "CONFIG"
	configDefaultPath = "config/secret-test.yaml"

	closeTimeout = 3
	readTimeout  = 10
	writeTimeout = 600
)

var (
	contextKey key
)

type key int

// start starts the server and waits for termination signal.
func start(ctx context.Context, router http.Handler) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.GetEnv(portEnvVar, portDefaultVal)),
		Handler:           router,
		ReadHeaderTimeout: readTimeout * time.Second,
		WriteTimeout:      writeTimeout * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			// adding server address to ctx handler functions receives
			return context.WithValue(ctx, contextKey, l.Addr().String())
		},
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("error listening for server")
		}
	}()
	log.Debug().Msg("server started")

	<-done
	log.Debug().Msg("server stopped")

	downCtx, cancel := context.WithTimeout(ctx, closeTimeout*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(downCtx); err != nil {
		log.Fatal().Err(err).Msg("error shuting server down")
	}
}

func getConfigOrPanic(version string) *config.Config {
	cnf, err := config.ReadFromFile(config.GetEnv(configEnvVar, configDefaultPath))
	if err != nil {
		log.Fatal().Err(err).Msg("error reading config")
	}
	cnf.Version = version
	return cnf
}

func initLogging(name, version, level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Warn().Err(err).Msgf("error parsing log level: %s", level)
	} else {
		zerolog.SetGlobalLevel(lvl)
	}

	log.Logger = zerolog.New(os.Stdout).With().
		Str("name", name).
		Str("version", version).
		Logger()
}
