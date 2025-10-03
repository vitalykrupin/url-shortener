-- Миграция для создания таблицы urls

CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    user_id TEXT,
    deleted_flag BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание индекса для ускорения поиска по alias
CREATE INDEX IF NOT EXISTS idx_urls_alias ON urls(alias);

-- Создание индекса для ускорения поиска по user_id
CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id);

-- Создание индекса для ускорения поиска по deleted_flag
CREATE INDEX IF NOT EXISTS idx_urls_deleted_flag ON urls(deleted_flag);