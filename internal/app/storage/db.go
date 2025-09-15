package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vitalykrupin/url-shortener/internal/app/middleware"
)

type DB struct {
	pool *pgxpool.Pool
}

var ErrDeleted = errors.New(`url deleted`)

func NewDB(DBDSN string) (*DB, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, DBDSN)
	if err != nil {
		log.Println("Can not connect to database")
		return nil, err
	}

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id serial PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			user_id TEXT,
			deleted_flag BOOLEAN NOT NULL DEFAULT FALSE
		);`)
	if err != nil {
		log.Println("Can not create table")
		return nil, err
	}

	return &DB{conn}, nil
}

func (d *DB) Add(ctx context.Context, batch map[Alias]OriginalURL) error {
	var query = `INSERT INTO urls (alias, url, user_id) VALUES (@alias, @url, @user_id)`
	b := pgx.Batch{}
	for alias, url := range batch {
		b.Queue(query, pgx.NamedArgs{
			"alias":   alias,
			"url":     url,
			"user_id": ctx.Value(middleware.UserIDKey),
		})
	}
	results := d.pool.SendBatch(ctx, &b)
	defer func(results pgx.BatchResults) {
		err := results.Close()
		if err != nil {
			log.Println("Can not close results")
		}
	}(results)

	for range batch {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("unable to insert row: %w", err)
		}
	}
	return nil
}

func (d *DB) GetAlias(ctx context.Context, url OriginalURL) (Alias, error) {
	var alias Alias
	err := d.pool.QueryRow(ctx, `SELECT alias FROM urls WHERE url = $1 AND user_id = $2 AND deleted_flag = false;`, url, ctx.Value(middleware.UserIDKey)).Scan(&alias)
	if err != nil {
		log.Println("Can not get alias from database", err)
		return "", err
	}
	return alias, nil
}

func (d *DB) GetURL(ctx context.Context, alias Alias) (OriginalURL, error) {
	var url OriginalURL
	var deletedFlag bool
	row := d.pool.QueryRow(ctx, `SELECT url, deleted_flag FROM urls WHERE alias = $1;`, alias)
	err := row.Scan(&url, &deletedFlag)
	if err != nil {
		log.Println("Can not get URL from database")
		return "", err
	}
	if deletedFlag {
		return "", ErrDeleted
	}
	return url, nil
}

func (d *DB) GetUserURLs(ctx context.Context, userID string) (aliasKeysMap AliasKeysMap, err error) {
	result := AliasKeysMap{}
	rows, err := d.pool.Query(ctx, `SELECT alias, url FROM urls WHERE user_id = $1;`, userID)
	if err != nil {
		log.Println("Can not get all user URLs from database")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var alias, originalURL string
		err := rows.Scan(&alias, &originalURL)
		if err != nil {
			return nil, errors.Join(errors.New("error scanning rows from rowset"), err)
		}
		result[Alias(alias)] = OriginalURL(originalURL)
	}
	return result, nil
}

func (d *DB) DeleteUserURLs(ctx context.Context, userID string, aliases []string) error {
	_, err := d.pool.Exec(ctx, `UPDATE urls SET deleted_flag = TRUE WHERE user_id = $1 AND alias = ANY($2);`, userID, aliases)
	if err != nil {
		log.Println("Can not delete URL from database", err)
		return err
	}
	return nil
}

func (d *DB) CloseStorage(ctx context.Context) error {
	d.pool.Close()
	return nil
}

func (d *DB) PingStorage(ctx context.Context) error {
	if err := d.pool.Ping(ctx); err != nil {
		log.Println("Can not ping database")
		return err
	}
	return nil
}
