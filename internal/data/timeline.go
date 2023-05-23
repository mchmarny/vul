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
	sqlTimelineList = `SELECT 
						x.imported,
						SUM(CASE WHEN source = 'grype' THEN 1 ELSE 0 END) grype,
						SUM(CASE WHEN source = 'trivy' THEN 1 ELSE 0 END) trivy,
						SUM(CASE WHEN source = 'snyk' THEN 1 ELSE 0 END) snyk
					  FROM (
						SELECT DISTINCT imported, source, package, version
						FROM vulns
						WHERE image = $1 
						AND imported >= $2
					  ) x
					  GROUP BY x.imported
					  ORDER BY 1`
)

func ListImageTimelines(ctx context.Context, pool *pgxpool.Pool, img, since string) ([]*vul.ImageTimeline, error) {
	if img == "" || since == "" {
		return nil, errors.New("empty image or since")
	}

	list := make([]*vul.ImageTimeline, 0)

	r := func(rows pgx.Rows) error {
		t := &vul.ImageTimeline{}
		if err := rows.Scan(
			&t.Date,
			&t.Grype,
			&t.Trivy,
			&t.Snyk,
		); err != nil {
			return errors.Wrapf(err, "failed to scan timeline row")
		}

		list = append(list, t)
		return nil
	}

	if err := mapRows(ctx, pool, r, sqlTimelineList, img, since); err != nil {
		return nil, errors.Wrap(err, "failed to map image version rows")
	}

	log.Info().Msgf("found %d timelines", len(list))

	return list, nil
}
