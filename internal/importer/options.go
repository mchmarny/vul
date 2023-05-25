package importer

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/internal/parser"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func ParseOptions(ctx context.Context, img, file, conn string) (*Options, error) {
	if img == "" {
		return nil, errors.New("image is required")
	}

	if file == "" {
		return nil, errors.New("file is required")
	}

	pool, err := data.GetPool(ctx, conn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create data pool")
	}

	o := &Options{
		Image: img,
		File:  file,
		Pool:  pool,
	}

	if !strings.Contains(o.Image, "@") {
		o.Image, err = parser.GetDigest(o.Image)
		if err != nil {
			return nil, errors.Wrap(err, "error getting digest")
		}
	}

	parts := strings.Split(o.Image, "@")
	o.ImageURI = parts[0]
	o.ImageDigest = parts[1]

	c, err := parser.GetContainer(o.File)
	if err != nil {
		return nil, errors.Wrap(err, "error getting container")
	}
	o.container = c
	o.Format = discoverFormat(c)

	if o.Format == FormatUnknown {
		return nil, errors.New("unknown format")
	}

	log.Debug().
		Str("image", o.Image).
		Str("digest", o.ImageDigest).
		Str("uri", o.ImageURI).
		Str("file", o.File).
		Str("format", o.Format.String()).
		Msg("parsed import options")

	return o, nil
}

// Options represents the input options.
type Options struct {
	Image  string
	File   string
	Format Format

	container *gabs.Container

	Pool        *pgxpool.Pool
	ImageURI    string
	ImageDigest string
}

// String returns a string representation of the options.
func (o *Options) String() string {
	return fmt.Sprintf("image: %s, file: %s, format: %s", o.Image, o.File, o.Format.String())
}
