package processor

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/mchmarny/vul/internal/importer"
	"github.com/mchmarny/vul/internal/parser"
	"github.com/mchmarny/vul/internal/scanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (h *Handler) processImage(ctx context.Context, imageURI string) error {
	if imageURI == "" {
		return errors.New("imageURI is empty")
	}

	// create temp dir
	dir := uuid.New().String()
	if err := os.Mkdir(dir, 0755); err != nil {
		return errors.Wrapf(err, "error creating temp dir: %s", dir)
	}
	defer os.RemoveAll(dir)

	// make sure image has digest
	imageDigest, err := parser.GetDigest(imageURI)
	if err != nil {
		return errors.Wrap(err, "error getting digest")
	}

	// scan image using digest into the temp dir
	reports, err := scanner.Scan(imageDigest, dir)
	if err != nil {
		return errors.Wrap(err, "error scanning image")
	}

	// import each resulting file
	for _, f := range reports {
		opt, err := importer.ParseOptions(imageURI, f, h.Pool)
		if err != nil {
			return errors.Wrapf(err, "error parsing options for %s with %s", f, imageURI)
		}
		if err := importer.Import(ctx, opt); err != nil {
			return errors.Wrapf(err, "error importing %s", opt)
		}
	}

	log.Debug().
		Str("image", imageURI).
		Int("reports", len(reports)).
		Msg("processed")

	return nil
}
