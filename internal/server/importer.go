package server

import (
	"context"
	"errors"

	"github.com/mchmarny/vul/internal/importer"
	"github.com/rs/zerolog/log"
)

// Import initiates the import process
// --version string   version of the application (default "v0.0.1")
// --image string     name of the image to import
// --file string      path to the file to import
func Import(version, image, file string) error {
	if image == "" || file == "" {
		err := errors.New("image, file, and conn are required")
		log.Error().
			Err(err).
			Str("image", image).
			Str("file", file).
			Msg("image and file are required")
		return err
	}

	cnf := getConfigOrPanic(version)

	initLogging("importer", version, cnf.Runtime.LogLevel)
	log.Debug().Msg("import initiated")

	ctx := context.Background()

	opt, err := importer.ParseOptions(ctx, cnf, image, file)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse options")
		return err
	}

	if err := importer.Import(ctx, opt); err != nil {
		log.Error().
			Err(err).
			Str("file", opt.File).
			Str("format", opt.Format.String()).
			Str("image", opt.Image).
			Str("image_uri", opt.ImageURI).
			Str("digest", opt.ImageDigest).
			Msg("failed to import")
		return err
	}

	log.Info().
		Str("file", opt.File).
		Str("format", opt.Format.String()).
		Str("image", opt.Image).
		Str("image_uri", opt.ImageURI).
		Str("digest", opt.ImageDigest).
		Msg("imported")

	return nil
}
