package data

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/pkg/query"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	sqlVersionList = `SELECT 
						digest,
						COUNT(distinct source) as sources,
						MIN(processed) as min_processed,
						MAX(processed) as max_processed,
						COUNT(distinct package) as packages
					  FROM vulns
					  WHERE image = $1
					  GROUP BY digest
					  ORDER BY 4 DESC`
)

func ListImageVersions(ctx context.Context, pool *pgxpool.Pool, imageURI string) ([]*query.ListImageVersionItem, error) {
	list := make([]*query.ListImageVersionItem, 0)

	r := func(rows pgx.Rows) error {
		q := &query.ListImageVersionItem{}
		if err := rows.Scan(&q.Digest, &q.SourceCount, &q.FirstReading, &q.LastReading, &q.PackageCount); err != nil {
			return errors.Wrapf(err, "failed to scan image version row")
		}
		list = append(list, q)
		return nil
	}

	if err := mapRows(ctx, pool, r, sqlVersionList, imageURI); err != nil {
		return nil, errors.Wrap(err, "failed to map image version rows")
	}

	log.Info().Msgf("found %d image versions", len(list))

	return list, nil
}
