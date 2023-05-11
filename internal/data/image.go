package data

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mchmarny/vul/pkg/vul"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	sqlImageList = `SELECT DISTINCT image FROM vulns ORDER BY 1`

	sqlSummary = `SELECT 
						COUNT(distinct image) images,
						COUNT(distinct digest) versions,
						COUNT(distinct source) sources,
						COUNT(distinct package) package,
						COUNT(*) exposure,
						SUM(CASE WHEN severity = 'negligible' THEN 1 ELSE 0 END) negligible,
						SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END) low,
						SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) medium,
						SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END) high,
						SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END) critical,
						SUM(CASE WHEN severity = 'unknown' THEN 1 ELSE 0 END) unknown,
						MIN(processed) fist_reading,
						MAX(processed) last_reading
					FROM vulns
					WHERE image = COALESCE($1, image)`
)

func ListImages(ctx context.Context, pool *pgxpool.Pool) ([]string, error) {
	list := make([]string, 0)

	r := func(rows pgx.Rows) error {
		var image string
		if err := rows.Scan(&image); err != nil {
			return errors.Wrapf(err, "failed to scan image row")
		}
		list = append(list, image)
		return nil
	}

	if err := mapRows(ctx, pool, r, sqlImageList); err != nil {
		return nil, errors.Wrap(err, "failed to map image rows")
	}

	log.Info().Msgf("found %d images", len(list))

	return list, nil
}

func GetSummary(ctx context.Context, pool *pgxpool.Pool, uri string) (*vul.SummaryItem, error) {
	img := &vul.SummaryItem{
		Image:    uri,
		Exposure: vul.Exposure{},
	}

	var arg sql.NullString
	if uri != "" {
		arg = sql.NullString{String: uri, Valid: true}
	}

	r := func(rows pgx.Row) error {
		if err := rows.Scan(
			&img.ImageCount,
			&img.VersionCount,
			&img.SourceCount,
			&img.PackageCount,
			&img.Exposure.Total,
			&img.Exposure.Negligible,
			&img.Exposure.Low,
			&img.Exposure.Medium,
			&img.Exposure.High,
			&img.Exposure.Critical,
			&img.Exposure.Unknown,
			&img.FirstReading,
			&img.LastReading); err != nil {
			return errors.Wrapf(err, "failed to scan image row")
		}
		return nil
	}

	if err := mapRow(ctx, pool, r, sqlSummary, arg); err != nil {
		return nil, errors.Wrap(err, "failed to map summary row")
	}

	return img, nil
}
