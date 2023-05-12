package data

import (
	"context"
	"fmt"

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
					  ORDER BY 1, 2`
)

func ListImageTimelines(ctx context.Context, pool *pgxpool.Pool, img, since string) ([]*vul.ImageTimeline, error) {
	if img == "" || since == "" {
		return nil, errors.New("empty image or since")
	}

	list := make([]*vul.ImageTimeline, 0)

	r := func(rows pgx.Rows) error {
		var date string
		var name string
		var total, negligible, low, medium, high, critical, unknown int

		if err := rows.Scan(
			&date,
			&name,
			&total,
			&negligible,
			&low,
			&medium,
			&high,
			&critical,
			&unknown,
		); err != nil {
			return errors.Wrapf(err, "failed to scan timeline row")
		}

		list = append(list, &vul.ImageTimeline{
			Date:  date,
			Name:  fmt.Sprintf("%s-total", name),
			Value: total,
		})

		list = append(list, &vul.ImageTimeline{
			Date:  date,
			Name:  fmt.Sprintf("%s-negligible", name),
			Value: negligible,
		})

		list = append(list, &vul.ImageTimeline{
			Date:  date,
			Name:  fmt.Sprintf("%s-low", name),
			Value: low,
		})

		list = append(list, &vul.ImageTimeline{
			Date:  date,
			Name:  fmt.Sprintf("%s-medium", name),
			Value: medium,
		})

		list = append(list, &vul.ImageTimeline{
			Date:  date,
			Name:  fmt.Sprintf("%s-high", name),
			Value: high,
		})

		list = append(list, &vul.ImageTimeline{
			Date:  date,
			Name:  fmt.Sprintf("%s-critical", name),
			Value: critical,
		})

		list = append(list, &vul.ImageTimeline{
			Date:  date,
			Name:  fmt.Sprintf("%s-unknown", name),
			Value: unknown,
		})
		return nil
	}

	if err := mapRows(ctx, pool, r, sqlTimelineList, img, since); err != nil {
		return nil, errors.Wrap(err, "failed to map image version rows")
	}

	log.Info().Msgf("found %d timelines", len(list))

	return list, nil
}
