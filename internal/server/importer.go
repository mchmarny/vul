package server

import (
	"context"

	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/internal/importer"
	"github.com/rs/zerolog/log"
)

// RunImport runs the import process.
func RunImport(version, image, file string) {
	cnf := getConfigOrPanic(version)

	if image == "" || file == "" {
		log.Fatal().
			Str("image", image).
			Str("file", file).
			Msg("image and file are required")
	}

	initLogging("cli", version, cnf.Runtime.LogLevel)
	log.Debug().Msg("import initiated")

	ctx := context.Background()

	pool, err := data.GetPool(ctx, cnf.Store)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create data pool")
	}

	opt, err := importer.ParseOptions(image, file, pool)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse options")
	}

	if err := importer.Import(ctx, opt); err != nil {
		log.Fatal().
			Err(err).
			Str("file", opt.File).
			Str("format", opt.Format.String()).
			Str("image", opt.Image).
			Str("image_uri", opt.ImageURI).
			Str("digest", opt.ImageDigest).
			Msg("failed to import")
	}

	log.Info().
		Str("file", opt.File).
		Str("format", opt.Format.String()).
		Str("image", opt.Image).
		Str("image_uri", opt.ImageURI).
		Str("digest", opt.ImageDigest).
		Msg("imported")
}
