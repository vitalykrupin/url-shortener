-- Миграция для создания таблицы users

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL UNIQUE
);

-- Создание индекса для ускорения поиска по login
CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);

-- Создание индекса для ускорения поиска по user_id
CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id);
