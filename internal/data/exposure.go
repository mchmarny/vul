package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	sqlExposureList = `SELECT 
						exposure,
						source,
						package,
						version,
						severity,
						MAX(score) as score,
						fixed
					  FROM vulns
					  WHERE image = $1 AND digest = $2
					  GROUP BY exposure, source, package, version, severity, fixed
					  ORDER BY 1, 2, 3, 4`
)

func ListImageVersionExposures(ctx context.Context, pool *pgxpool.Pool, imageURI, digest string) (map[string][]*vul.ListDigestExposureItem, error) {
	list := make(map[string][]*vul.ListDigestExposureItem)

	r := func(rows pgx.Rows) error {
		q := &vul.ListDigestExposureItem{}
		var e string
		if err := rows.Scan(
			&e,
			&q.Source,
			&q.Package,
			&q.Version,
			&q.Severity,
			&q.Score,
			&q.Fixed); err != nil {
			return errors.Wrapf(err, "failed to scan image version row")
		}
		if _, ok := list[e]; !ok {
			list[e] = make([]*vul.ListDigestExposureItem, 0)
		}

		list[e] = append(list[e], q)
		return nil
	}

	if err := mapRows(ctx, pool, r, sqlExposureList, imageURI, digest); err != nil {
		return nil, errors.Wrap(err, "failed to map image version exposure rows")
	}

	log.Info().Msgf("found %d image versions", len(list))

	return list, nil
}
