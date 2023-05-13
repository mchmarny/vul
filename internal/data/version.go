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
	sqlVersionList = `SELECT 
						digest,
						MAX(processed) processed
					  FROM vulns
					  WHERE image = $1
					  GROUP BY digest
					  ORDER BY 2 DESC`
)

func ListImageVersions(ctx context.Context, pool *pgxpool.Pool, imageURI string) ([]*vul.ImageVersion, error) {
	list := make([]*vul.ImageVersion, 0)

	r := func(rows pgx.Rows) error {
		q := &vul.ImageVersion{
			Image: imageURI,
		}
		if err := rows.Scan(
			&q.Digest,
			&q.Processed); err != nil {
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
