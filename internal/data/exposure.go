package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/pkg/errors"
)

var (
	sqlExposureList = `SELECT 
						package,
						version,
						source,
						exposure,
						severity,
						MAX(score) as score,
						fixed
					  FROM vulns
					  WHERE image = $1 AND digest = $2
					  GROUP BY package, version, source, exposure, severity, fixed
					  ORDER BY 1, 2, 3, 4`
)

func ListImageVersionExposures(ctx context.Context, pool *pgxpool.Pool, imageURI, digest string) (*vul.ImageDigestExposures, error) {
	if imageURI == "" || digest == "" {
		return nil, errors.New("image and digest are required")
	}

	m := &vul.ImageDigestExposures{
		Image:    imageURI,
		Digest:   digest,
		Packages: make(map[string]*vul.PackageExposures),
	}

	r := func(rows pgx.Rows) error {
		var pkg, ver, src, exp, sev string
		var score float64
		var fixed bool

		if err := rows.Scan(
			&pkg,
			&ver,
			&src,
			&exp,
			&sev,
			&score,
			&fixed); err != nil {
			return errors.Wrapf(err, "failed to scan image version exposure row")
		}

		if fixed {
			m.FixedCount++
		}

		if _, ok := m.Packages[pkg]; !ok {
			m.Packages[pkg] = &vul.PackageExposures{
				Versions: make(map[string]*vul.PackageVersionExposures),
			}
		}

		if _, ok := m.Packages[pkg].Versions[ver]; !ok {
			m.Packages[pkg].Versions[ver] = &vul.PackageVersionExposures{
				Sources: make(map[string]*vul.SourceExposures),
			}
		}

		if _, ok := m.Packages[pkg].Versions[ver].Sources[src]; !ok {
			m.Packages[pkg].Versions[ver].Sources[src] = &vul.SourceExposures{
				Exposures: make(map[string]*vul.Exposures),
			}
		}

		if _, ok := m.Packages[pkg].Versions[ver].Sources[src].Exposures[exp]; !ok {
			m.Packages[pkg].Versions[ver].Sources[src].Exposures[exp] = &vul.Exposures{
				Severity: sev,
				Score:    score,
				Fixed:    fixed,
			}
		}

		return nil
	}

	if err := mapRows(ctx, pool, r, sqlExposureList, imageURI, digest); err != nil {
		return nil, errors.Wrap(err, "failed to map image version exposure rows")
	}

	return m, nil
}
