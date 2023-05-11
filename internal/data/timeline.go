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
						x.source,
						COUNT(*) total,
						SUM(CASE WHEN x.severity = 'negligible' THEN 1 ELSE 0 END) negligible,
						SUM(CASE WHEN x.severity = 'low' THEN 1 ELSE 0 END) low,
						SUM(CASE WHEN x.severity = 'medium' THEN 1 ELSE 0 END) medium,
						SUM(CASE WHEN x.severity = 'high' THEN 1 ELSE 0 END) high,
						SUM(CASE WHEN x.severity = 'critical' THEN 1 ELSE 0 END) critical,
						SUM(CASE WHEN x.severity = 'unknown' THEN 1 ELSE 0 END) unknown
					  FROM (
						SELECT DISTINCT imported, source, exposure, severity, package, version
						FROM vulns
						WHERE image = $1 
						AND imported >= $2
					  ) x
					  GROUP BY x.imported, x.source
					  ORDER BY 1 DESC, 2`
)

func ListImageTimelines(ctx context.Context, pool *pgxpool.Pool, img, since string) (map[string]*vul.ListImageTimelineItem, error) {
	if img == "" || since == "" {
		return nil, errors.New("empty image or since")
	}

	m := make(map[string]*vul.ListImageTimelineItem)

	r := func(rows pgx.Rows) error {
		q := &vul.ListImageSourceTimelineItem{}
		var day string
		var src string
		if err := rows.Scan(
			&day,
			&src,
			&q.Total,
			&q.Negligible,
			&q.Low,
			&q.Medium,
			&q.High,
			&q.Critical,
			&q.Unknown,
		); err != nil {
			return errors.Wrapf(err, "failed to scan timeline row")
		}

		if _, ok := m[day]; !ok {
			m[day] = &vul.ListImageTimelineItem{
				Sources: make(map[string]*vul.ListImageSourceTimelineItem),
			}
		}

		m[day].Sources[src] = q
		return nil
	}

	if err := mapRows(ctx, pool, r, sqlTimelineList, img, since); err != nil {
		return nil, errors.Wrap(err, "failed to map image version rows")
	}

	log.Info().Msgf("found %d timelines", len(m))

	return m, nil
}
