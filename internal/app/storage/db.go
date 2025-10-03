// Package storage предоставляет реализацию хранилища данных в PostgreSQL
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

// DB реализация хранилища данных в PostgreSQL
type DB struct {
	// pool пул соединений с базой данных
	pool *pgxpool.Pool
}

// ErrDeleted ошибка, возникающая при попытке получить удаленный URL
var ErrDeleted = errors.New(`url deleted`)

// NewDB создает новое подключение к базе данных PostgreSQL
// DBDSN строка подключения к базе данных
// Возвращает указатель на DB и ошибку, если подключение не удалось
func NewDB(DBDSN string) (*DB, error) {
	ctx := context.Background()

	// Создание пула соединений
	conn, err := pgxpool.New(ctx, DBDSN)
	if err != nil {
		log.Println("Can not connect to database")
		return nil, err
	}

	// Создание таблицы urls, если она не существует
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

// Add добавляет новые URL в базу данных
// ctx контекст запроса
// batch карта alias -> OriginalURL для добавления
// Возвращает ошибку, если добавление не удалось
func (d *DB) Add(ctx context.Context, batch map[Alias]OriginalURL) error {
	if len(batch) == 0 {
		return nil
	}

	var query = `INSERT INTO urls (alias, url, user_id) VALUES (@alias, @url, @user_id)`
	b := &pgx.Batch{}
	for alias, url := range batch {
		b.Queue(query, pgx.NamedArgs{
			"alias":   alias,
			"url":     url,
			"user_id": ctx.Value(middleware.UserIDKey),
		})
	}

	results := d.pool.SendBatch(ctx, b)
	defer func() {
		if err := results.Close(); err != nil {
			log.Printf("Failed to close batch results: %v", err)
		}
	}()

	for i := 0; i < len(batch); i++ {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("failed to execute batch query #%d: %w", i, err)
		}
	}

	return nil
}

// GetAlias получает alias для заданного URL
// ctx контекст запроса
// url оригинальный URL
// Возвращает alias и ошибку, если получение не удалось
func (d *DB) GetAlias(ctx context.Context, url OriginalURL) (Alias, error) {
	userID := ctx.Value(middleware.UserIDKey)
	if userID == nil || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}

	var alias Alias
	err := d.pool.QueryRow(ctx, `SELECT alias FROM urls WHERE url = $1 AND user_id = $2 AND deleted_flag = false;`, url, userID).Scan(&alias)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("alias not found for URL: %s", url)
		}
		log.Printf("Failed to get alias from database: %v", err)
		return "", fmt.Errorf("database error: %w", err)
	}
	return alias, nil
}

// GetURL получает оригинальный URL по alias
// ctx контекст запроса
// alias короткий alias URL
// Возвращает оригинальный URL и ошибку, если получение не удалось
func (d *DB) GetURL(ctx context.Context, alias Alias) (OriginalURL, error) {
	var url OriginalURL
	var deletedFlag bool
	row := d.pool.QueryRow(ctx, `SELECT url, deleted_flag FROM urls WHERE alias = $1;`, alias)
	err := row.Scan(&url, &deletedFlag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("URL not found for alias: %s", alias)
		}
		log.Printf("Failed to get URL from database: %v", err)
		return "", fmt.Errorf("database error: %w", err)
	}
	if deletedFlag {
		return "", ErrDeleted
	}
	return url, nil
}

// GetUserURLs получает все URL пользователя
// ctx контекст запроса
// userID идентификатор пользователя
// Возвращает карту alias -> OriginalURL и ошибку, если получение не удалось
func (d *DB) GetUserURLs(ctx context.Context, userID string) (AliasKeysMap, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	result := make(AliasKeysMap)
	rows, err := d.pool.Query(ctx, `SELECT alias, url FROM urls WHERE user_id = $1;`, userID)
	if err != nil {
		log.Printf("Failed to query user URLs from database: %v", err)
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer func() {
		rows.Close()
	}()

	for rows.Next() {
		var alias, originalURL string
		if err := rows.Scan(&alias, &originalURL); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result[Alias(alias)] = OriginalURL(originalURL)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}

// DeleteUserURLs помечает URL пользователя как удаленные
// ctx контекст запроса
// userID идентификатор пользователя
// aliases список alias для удаления
// Возвращает ошибку, если удаление не удалось
func (d *DB) DeleteUserURLs(ctx context.Context, userID string, aliases []string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	if len(aliases) == 0 {
		return nil
	}

	_, err := d.pool.Exec(ctx, `UPDATE urls SET deleted_flag = TRUE WHERE user_id = $1 AND alias = ANY($2);`, userID, aliases)
	if err != nil {
		log.Printf("Failed to delete URLs from database: %v", err)
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

// CloseStorage закрывает подключение к базе данных
// ctx контекст запроса
// Возвращает ошибку, если закрытие не удалось
func (d *DB) CloseStorage(ctx context.Context) error {
	if d.pool != nil {
		d.pool.Close()
	}
	return nil
}

// PingStorage проверяет подключение к базе данных
// ctx контекст запроса
// Возвращает ошибку, если подключение не удалось
func (d *DB) PingStorage(ctx context.Context) error {
	if d.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	if err := d.pool.Ping(ctx); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}
