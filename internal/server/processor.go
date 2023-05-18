package server

import (
	"context"

	"github.com/mchmarny/vul/internal/processor"
	"github.com/rs/zerolog/log"
)

func RunProcessor(version string) {
	cnf := getConfigOrPanic(version)

	initLogging("processor", version, cnf.Runtime.LogLevel)
	log.Info().
		Str("name", cnf.Name).
		Str("version", version).
		Msg("starting processor server")

	ctx := context.Background()
	h, err := processor.New(ctx, cnf)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating processor handler")
	}
	defer h.Close()

	start(ctx, h.Router)
}
