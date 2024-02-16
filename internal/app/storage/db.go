package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/vitalykrupin/url-shortener.git/cmd/shortener/config"
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
	// defer conn.Close(context.Background())

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS urls (id serial PRIMARY KEY, alias TEXT NOT NULL UNIQUE, url TEXT NOT NULL)`)
	if err != nil {
		log.Println("Can not create table")
		return nil, err
	}

	return &DB{conn}, nil
}

func (d *DB) Add(ctx context.Context, alias string, url string) error {
	_, err := d.conn.Exec(ctx, `INSERT INTO urls (alias, url) VALUES ($1, $2);`, alias, url)
	return err
}

func (d *DB) GetURL(ctx context.Context, alias string) (string, error) {
	var url string
	err := d.conn.QueryRow(ctx, `SELECT url FROM urls WHERE alias = $1;`, alias).Scan(&url)
	if err != nil {
		log.Println("Can not get URL from database")
		return "", err
	}
	return url, nil
}

func (d *DB) GetAlias(ctx context.Context, url string) (string, error) {
	var alias string
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
