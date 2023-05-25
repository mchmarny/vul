package data

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func GetPool(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to database: %s", uri)
	}

	conn, err := getDBConn(ctx, pool)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get connection from pool")
	}
	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		return nil, errors.Wrapf(err, "failed to ping database")
	}

	return pool, nil
}

func getDBConn(ctx context.Context, pool *pgxpool.Pool) (*pgxpool.Conn, error) {
	if pool == nil {
		return nil, errors.New("nil pool")
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to acquire connection")
	}

	return conn, nil
}

type rowsMapper func(rows pgx.Rows) error
type rowMapper func(rows pgx.Row) error

func mapRows(ctx context.Context, p *pgxpool.Pool, m rowsMapper, q string, args ...any) error {
	conn, err := getDBConn(ctx, p)
	if err != nil {
		return errors.Wrapf(err, "failed to get db conn")
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, q, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Debug().
			Err(err).
			Str("query", sqlImageList).
			Interface("args", args).
			Msg("error executing query")
		return errors.Wrapf(err, "failed to execute select statement")
	}
	defer rows.Close()

	for rows.Next() {
		if err := m(rows); err != nil {
			return errors.Wrapf(err, "failed to map row")
		}
	}

	return nil
}

func mapRow(ctx context.Context, p *pgxpool.Pool, m rowMapper, q string, args ...any) error {
	conn, err := getDBConn(ctx, p)
	if err != nil {
		return errors.Wrapf(err, "failed to get db conn")
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, q, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Debug().
			Err(err).
			Str("query", sqlImageList).
			Interface("args", args).
			Msg("error executing query")
		return errors.Wrapf(err, "failed to execute select statement")
	}

	if err := m(row); err != nil {
		return errors.Wrapf(err, "failed to map row")
	}

	return nil
}
