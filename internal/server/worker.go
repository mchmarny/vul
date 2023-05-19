package server

import (
	"context"

	"github.com/mchmarny/vul/internal/worker"
	"github.com/rs/zerolog/log"
)

func RunWorker(version string) {
	cnf := getConfigOrPanic(version)

	initLogging("worker", version, cnf.Runtime.LogLevel)
	log.Info().
		Str("name", cnf.Name).
		Str("version", version).
		Msg("starting worker server")

	ctx := context.Background()
	h, err := worker.New(ctx, cnf)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating worker handler")
	}
	defer h.Close()

	start(ctx, h.Router)
}
