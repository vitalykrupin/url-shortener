package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
	"github.com/vitalykrupin/url-shortener.git/internal/app/middleware"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB(ctx context.Context, cfg *config.Config) (Storage, error) {
	if cfg.DBDSN == "" {
		log.Println("No DBDSN provided")
		return nil, fmt.Errorf("no DBDSN provided")
	}
	conn, err := pgx.Connect(ctx, cfg.DBDSN)
	if err != nil {
		log.Println("Can not connect to database")
		return nil, err
	}

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id serial PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			user_id TEXT
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
	results := d.conn.SendBatch(ctx, &b)
	defer results.Close()

	for range batch {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("unable to insert row: %w", err)
		}
	}
	return nil
}

func (d *DB) GetURL(ctx context.Context, alias Alias) (OriginalURL, error) {
	var url OriginalURL
	err := d.conn.QueryRow(ctx, `SELECT url FROM urls WHERE alias = $1;`, alias).Scan(&url)
	if err != nil {
		log.Println("Can not get URL from database")
		return "", err
	}
	return url, nil
}

func (d *DB) GetUserURLs(ctx context.Context, userID string) (*aliasKeysMap, error) {
	result := aliasKeysMap{}
	rows, err := d.conn.Query(ctx, `SELECT alias, url FROM urls WHERE user_id = $1;`, userID)
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
	return &result, nil
}

func (d *DB) GetAlias(ctx context.Context, url OriginalURL) (Alias, error) {
	var alias Alias
	err := d.conn.QueryRow(ctx, `SELECT alias FROM urls WHERE url = $1;`, url).Scan(&alias)
	if err != nil {
		log.Println("Can not get alias from database")
		return "", err
	}
	return alias, nil
}

func (d *DB) CloseStorage(ctx context.Context) error {
	if err := d.conn.Close(ctx); err != nil {
		return fmt.Errorf("error database closing: %w", err)
	}
	return nil
}

func (d *DB) PingStorage(ctx context.Context) error {
	if err := d.conn.Ping(ctx); err != nil {
		log.Println("Can not ping database")
		return err
	}
	return nil
}
