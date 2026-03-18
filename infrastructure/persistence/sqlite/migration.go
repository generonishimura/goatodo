package sqlite

import "database/sql"

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			status TEXT NOT NULL DEFAULT 'todo',
			priority INTEGER NOT NULL DEFAULT 0,
			due_date TEXT,
			created_at TEXT NOT NULL,
			completed_at TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
		CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	`)
	return err
}
