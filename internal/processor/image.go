package processor

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/mchmarny/vul/internal/scanner"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func processImage(ctx context.Context, imageURI string) error {
	if imageURI == "" {
		return errors.New("imageURI is empty")
	}

	dir := uuid.New().String()

	if err := os.Mkdir(dir, 0755); err != nil {
		return errors.Wrapf(err, "error creating temp dir: %s", dir)
	}
	defer os.RemoveAll(dir)

	filesExpected := scanner.Scan(imageURI, dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		return errors.Wrapf(err, "error reading dir: %s", dir)
	}

	if len(files) != filesExpected {
		return errors.Errorf("expected %d files, got %d, see logs for details", filesExpected, len(files))
	}

	for _, f := range files {
		log.Debug().Str("file", f.Name()).Msg("processing file")
	}

	return nil
}
