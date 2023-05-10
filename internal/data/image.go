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
	sqlImageList = `SELECT 
						image, 
						COUNT(distinct digest) versions,
						MIN(processed) fist_reading,
						MAX(processed) last_reading
					FROM vulns 
					GROUP BY image`
)

func ListImages(ctx context.Context, pool *pgxpool.Pool) ([]*query.ListImageItem, error) {
	conn, err := getDBConn(ctx, pool)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get db conn")
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, sqlImageList)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Debug().
			Err(err).
			Str("query", sqlImageList).
			Msg("error executing query")
		return nil, errors.Wrapf(err, "failed to execute select statement")
	}
	defer rows.Close()

	list := make([]*query.ListImageItem, 0)

	for rows.Next() {
		q := &query.ListImageItem{}
		if err := rows.Scan(&q.Image, &q.VersionCount, &q.FirstReading, &q.LastReading); err != nil {
			return nil, errors.Wrapf(err, "failed to scan image row")
		}
		list = append(list, q)
	}

	log.Info().Msgf("found %d records", len(list))

	return list, nil
}
