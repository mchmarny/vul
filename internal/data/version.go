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
						source,
						COUNT(distinct package) as packages,
						MIN(processed) as min_processed,
						MAX(processed) as max_processed
					  FROM vulns
					  WHERE image = $1
					  GROUP BY digest, source
					  ORDER BY 1, 2, 3`
)

func ListImageVersions(ctx context.Context, pool *pgxpool.Pool, imageURI string) (map[string][]*query.ListImageSourceItem, error) {
	list := make(map[string][]*query.ListImageSourceItem)

	r := func(rows pgx.Rows) error {
		q := &query.ListImageSourceItem{}
		var d string
		if err := rows.Scan(
			&d,
			&q.Source,
			&q.PackageCount,
			&q.FirstReading,
			&q.LastReading); err != nil {
			return errors.Wrapf(err, "failed to scan image version row")
		}
		if _, ok := list[d]; !ok {
			list[d] = make([]*query.ListImageSourceItem, 0)
		}

		list[d] = append(list[d], q)
		return nil
	}

	if err := mapRows(ctx, pool, r, sqlVersionList, imageURI); err != nil {
		return nil, errors.Wrap(err, "failed to map image version rows")
	}

	log.Info().Msgf("found %d image versions", len(list))

	return list, nil
}
