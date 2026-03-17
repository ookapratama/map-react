-- Rollback: 000001_create_tasks_table
-- Description: Drop the tasks table and its indexes

DROP INDEX IF EXISTS idx_tasks_search;
DROP INDEX IF EXISTS idx_tasks_status;
DROP TABLE IF EXISTS tasks;
