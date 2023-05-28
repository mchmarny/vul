package importer

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/internal/config"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/internal/metric"
	"github.com/mchmarny/vul/internal/parser"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Format represents the source format.
func ParseOptions(ctx context.Context, cnf *config.Config, img, file string) (*Options, error) {
	if cnf == nil {
		return nil, errors.New("config is required")
	}

	if img == "" {
		return nil, errors.New("image is required")
	}

	if file == "" {
		return nil, errors.New("file is required")
	}

	// example: postgres://postgres:***@10.126.160.7:5432/vul
	conn := fmt.Sprintf("%s://%s:%s@%s:%d/%s", cnf.Store.Type, cnf.Store.User, cnf.Store.Password, cnf.Store.Host, cnf.Store.Port, cnf.Store.DB)

	pool, err := data.GetPool(ctx, conn)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("type", cnf.Store.Type).
			Str("user", cnf.Store.User).
			Str("host", cnf.Store.Host).
			Str("db", cnf.Store.DB).
			Msg("failed to create data pool")
	}

	rec, err := metric.New(cnf.ProjectID, cnf.Name, cnf.Version, cnf.Runtime.SendMetrics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create metric service")
	}

	o := &Options{
		Image:    img,
		File:     file,
		Config:   cnf,
		Pool:     pool,
		Recorder: rec,
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

	Config      *config.Config
	Pool        *pgxpool.Pool
	Recorder    metric.Recorder
	ImageURI    string
	ImageDigest string
}

// String returns a string representation of the options.
func (o *Options) String() string {
	return fmt.Sprintf("image: %s, file: %s, format: %s", o.Image, o.File, o.Format.String())
}
