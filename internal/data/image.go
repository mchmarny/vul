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
					COUNT(DISTINCT image) images,
					COUNT(DISTINCT digest) versions,
					COUNT(DISTINCT source) sources,
					COUNT(DISTINCT package) package,
					COUNT(exposure) exposure,
					SUM(CASE WHEN severity = 'negligible' THEN 1 ELSE 0 END) negligible,
					SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END) low,
					SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) medium,
					SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END) high,
					SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END) critical,
					SUM(CASE WHEN severity = 'unknown' THEN 1 ELSE 0 END) unknown,
					MIN(min_processed) min_processed,
					MAX(max_processed) max_processed
				FROM (
					SELECT 
						image, 
						digest, 
						source, 
						package, 
						exposure, 
						severity, 
						MIN(processed) min_processed,
						MAX(processed) max_processed
					FROM vulns
					WHERE image = COALESCE($1, image) AND digest = COALESCE($2, digest)
					GROUP BY image, digest, source, package, exposure, severity
				) x
				`
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

func GetSummary(ctx context.Context, pool *pgxpool.Pool, img, dig string) (*vul.SummaryItem, error) {
	s := &vul.SummaryItem{
		Image:    img,
		Exposure: vul.Exposure{},
	}

	var imgArg sql.NullString
	if img != "" {
		imgArg = sql.NullString{String: img, Valid: true}
	}

	var digArg sql.NullString
	if dig != "" {
		digArg = sql.NullString{String: dig, Valid: true}
	}

	r := func(rows pgx.Row) error {
		if err := rows.Scan(
			&s.ImageCount,
			&s.VersionCount,
			&s.SourceCount,
			&s.PackageCount,
			&s.TotalExposures,
			&s.Exposure.Negligible,
			&s.Exposure.Low,
			&s.Exposure.Medium,
			&s.Exposure.High,
			&s.Exposure.Critical,
			&s.Exposure.Unknown,
			&s.FirstReading,
			&s.LastReading); err != nil {
			return errors.Wrapf(err, "failed to scan image row")
		}
		return nil
	}

	if err := mapRow(ctx, pool, r, sqlSummary, imgArg, digArg); err != nil {
		return nil, errors.Wrap(err, "failed to map summary row")
	}

	s.DaysCount = int(s.LastReading.Sub(s.FirstReading).Hours() / 24)
	if s.SourceCount > 0 {
		s.AvgExposure = float64(s.TotalExposures / s.SourceCount)
	}

	return s, nil
}
