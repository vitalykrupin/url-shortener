package storage

import (
	"context"
	"database/sql"
	"log"
)

type DB struct {
	conn *sql.DB
}

func NewDB(conn *sql.DB) *DB {
	return &DB{conn: conn}
}

func (d *DB) BootstrapDB(ctx context.Context) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY, alias TEXT UNIQUE, original_url TEXT)")
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *DB) Add(ctx context.Context, alias string, url string) error {
	_, err := d.conn.ExecContext(ctx, "INSERT INTO urls (alias, url) VALUES ($1, $2);", alias, url)
	return err
}

func (d *DB) GetURL(ctx context.Context, alias string) (string, error) {
	var url string
	err := d.conn.QueryRowContext(ctx, "SELECT url FROM urls WHERE alias = $1;", alias).Scan(&url)
	if err != nil {
		log.Println("Can not get URL from database")
		return "", err
	}
	return url, nil
}

func (d *DB) GetAlias(ctx context.Context, url string) (string, error) {
	var alias string
	err := d.conn.QueryRowContext(ctx, "SELECT alias FROM urls WHERE url = $1;", url).Scan(&alias)
	if err != nil {
		log.Println("Can not get alias from database")
		return "", err
	}

	return alias, nil
}
