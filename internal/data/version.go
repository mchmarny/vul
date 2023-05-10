package data

import (
	"context"
	"database/sql"

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
	conn, err := getDBConn(ctx, pool)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get db conn")
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, sqlVersionList, imageURI)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Debug().
			Err(err).
			Str("query", sqlImageList).
			Interface("params", imageURI).
			Msg("error executing query")
		return nil, errors.Wrapf(err, "failed to execute select statement")
	}
	defer rows.Close()

	list := make([]*query.ListImageVersionItem, 0)

	for rows.Next() {
		q := &query.ListImageVersionItem{}
		if err := rows.Scan(&q.Digest, &q.SourceCount, &q.FirstReading, &q.LastReading, &q.PackageCount); err != nil {
			return nil, errors.Wrapf(err, "failed to scan image version row")
		}
		list = append(list, q)
	}

	log.Info().Msgf("found %d records", len(list))

	return list, nil
}
