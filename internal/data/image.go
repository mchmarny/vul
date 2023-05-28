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
						COUNT(DISTINCT package) packages, 
						COUNT(exposure) exposures,
						COUNT(DISTINCT exposure) unique_exposures,
						COUNT(DISTINCT exposure) FILTER (where fixed = true) fixed,
						COUNT(severity) FILTER (where severity = 'negligible') negligible,
						COUNT(severity) FILTER (where severity = 'low') low,
						COUNT(severity) FILTER (where severity = 'medium') medium,
						COUNT(severity) FILTER (where severity = 'high') high,
						COUNT(severity) FILTER (where severity = 'critical') critical,
						COUNT(severity) FILTER (where severity = 'unknown') unknown,
						MAX(processed) last_reading
					FROM vulns
					WHERE imported = (SELECT MAX(imported) FROM vulns)
					AND image = COALESCE($1, image) 
					AND digest = COALESCE($2, digest)
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
			&s.UniqueExposures,
			&s.FixedCount,
			&s.Exposure.Negligible,
			&s.Exposure.Low,
			&s.Exposure.Medium,
			&s.Exposure.High,
			&s.Exposure.Critical,
			&s.Exposure.Unknown,
			&s.LastReading); err != nil {
			return errors.Wrapf(err, "failed to scan image row")
		}
		return nil
	}

	if err := mapRow(ctx, pool, r, sqlSummary, imgArg, digArg); err != nil {
		return nil, errors.Wrap(err, "failed to map summary row")
	}

	return s, nil
}
