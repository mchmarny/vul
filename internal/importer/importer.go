package importer

import (
	"context"

	"github.com/Jeffail/gabs/v2"
	"github.com/mchmarny/vul/internal/converter/grype"
	"github.com/mchmarny/vul/internal/converter/snyk"
	"github.com/mchmarny/vul/internal/converter/trivy"
	"github.com/mchmarny/vul/internal/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func Import(ctx context.Context, opt *Options) error {
	if opt == nil {
		return errors.New("options required")
	}

	m, err := getMapper(opt.Format)
	if err != nil {
		return errors.Wrap(err, "error getting converter")
	}

	list, err := m(opt.container)
	if err != nil {
		return errors.Wrap(err, "error converting source")
	}

	if list == nil {
		return errors.New("expected non-nil array of vulnerabilities")
	}

	vuls := data.DecorateVulnerabilities(unique(list), opt.ImageURI, opt.ImageDigest, opt.Format.String())
	log.Debug().Int("vulnerabilities", len(vuls)).Msg("found")

	if err := data.Import(ctx, opt.Pool, vuls); err != nil {
		return errors.Wrap(err, "error importing vulnerabilities")
	}

	return nil
}

// VulnerabilityMapper is a function that converts a source to a list of common vulnerability types.
type VulnerabilityMapper func(c *gabs.Container) ([]*data.Vulnerability, error)

// GetMapper returns a vulnerability converter for the given source format.
func getMapper(format Format) (VulnerabilityMapper, error) {
	switch format {
	case FormatSnykJSON:
		return snyk.Convert, nil
	case FormatTrivyJSON:
		return trivy.Convert, nil
	case FormatGrypeJSON:
		return grype.Convert, nil
	default:
		return nil, errors.Errorf("unimplemented conversion format: %s", format)
	}
}

// Hashible is an interface for objects that can be hashed.
type Hashible interface {
	GetID() string
}

func unique[T Hashible](list []T) []T {
	seen := map[string]bool{}
	result := make([]T, 0)
	for _, item := range list {
		h := item.GetID()
		if !seen[h] {
			seen[h] = true
			result = append(result, item)
		}
	}
	return result
}
