-- Откат миграции для удаления таблицы users

-- Удаление индексов
DROP INDEX IF EXISTS idx_users_login;
DROP INDEX IF EXISTS idx_users_user_id;

-- Удаление таблицы
DROP TABLE IF EXISTS users;
