package server

import (
	"context"

	"github.com/mchmarny/vul/internal/handler"
	"github.com/rs/zerolog/log"
)

// RunApp starts the server and waits for termination signal.
func RunApp(version string) {
	cnf := getConfigOrPanic(version)

	initLogging("app", version, cnf.Runtime.LogLevel)
	log.Info().
		Str("name", cnf.Name).
		Str("version", version).
		Msg("starting app server")

	ctx := context.Background()
	h, err := handler.New(ctx, cnf)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating app handler")
	}
	defer h.Close()

	start(ctx, h.Router)
}
