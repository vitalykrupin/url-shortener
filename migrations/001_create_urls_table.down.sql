-- Откат миграции для удаления таблицы urls

-- Удаление индексов
DROP INDEX IF EXISTS idx_urls_alias;
DROP INDEX IF EXISTS idx_urls_user_id;
DROP INDEX IF EXISTS idx_urls_deleted_flag;

-- Удаление таблицы
DROP TABLE IF EXISTS urls;