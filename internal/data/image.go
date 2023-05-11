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
	sqlImageList = `SELECT 
						image, 
						COUNT(distinct digest) versions,
						COUNT(distinct source) sources,
						MIN(processed) fist_reading,
						MAX(processed) last_reading
					FROM vulns 
					GROUP BY image`
)

func ListImages(ctx context.Context, pool *pgxpool.Pool) ([]*vul.ListImageItem, error) {
	list := make([]*vul.ListImageItem, 0)

	r := func(rows pgx.Rows) error {
		q := &vul.ListImageItem{}
		if err := rows.Scan(
			&q.Image,
			&q.VersionCount,
			&q.SourceCount,
			&q.FirstReading,
			&q.LastReading); err != nil {
			return errors.Wrapf(err, "failed to scan image row")
		}
		list = append(list, q)
		return nil
	}

	if err := mapRows(ctx, pool, r, sqlImageList); err != nil {
		return nil, errors.Wrap(err, "failed to map image rows")
	}

	log.Info().Msgf("found %d images", len(list))

	return list, nil
}
