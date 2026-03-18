package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/i-nishimura/goatodo/infrastructure/persistence/sqlite"
	"github.com/i-nishimura/goatodo/presenter"
	_ "modernc.org/sqlite"
)

type App struct {
	ctx         context.Context
	db          *sql.DB
	TaskHandler *presenter.TaskHandler
}

func NewApp() *App {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := sqlite.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	repo := sqlite.NewTaskRepository(db)
	taskHandler := presenter.NewTaskHandler(repo)

	return &App{
		db:          db,
		TaskHandler: taskHandler,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

func getDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "goatodo.db"
	}
	dir := filepath.Join(homeDir, ".goatodo")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "goatodo.db")
}
