package sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/i-nishimura/goatodo/domain/task"
	sqliterepo "github.com/i-nishimura/goatodo/infrastructure/persistence/sqlite"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := sqliterepo.Migrate(db); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestSQLiteTaskRepository_SaveAndFindByID(t *testing.T) {
	t.Run("saves and retrieves a task", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewTaskRepository(db)

		result := task.NewTask("Buy milk")
		if result.IsErr() {
			t.Fatalf("setup failed: %s", result.Error())
		}
		tk := result.Value()

		saveResult := repo.Save(tk)
		if saveResult.IsErr() {
			t.Fatalf("save failed: %s", saveResult.Error())
		}

		found := repo.FindByID(tk.ID())
		if found.IsErr() {
			t.Fatalf("find failed: %s", found.Error())
		}

		if found.Value().ID() != tk.ID() {
			t.Errorf("expected ID %s, got %s", tk.ID(), found.Value().ID())
		}
		if found.Value().Title() != "Buy milk" {
			t.Errorf("expected title 'Buy milk', got '%s'", found.Value().Title())
		}
		if found.Value().Status() != task.StatusTodo {
			t.Errorf("expected status todo, got %s", found.Value().Status())
		}
	})
}

func TestSQLiteTaskRepository_FindAll(t *testing.T) {
	t.Run("returns all saved tasks", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewTaskRepository(db)

		r1 := task.NewTask("Task 1")
		r2 := task.NewTask("Task 2")
		repo.Save(r1.Value())
		repo.Save(r2.Value())

		all := repo.FindAll()
		if all.IsErr() {
			t.Fatalf("find all failed: %s", all.Error())
		}
		if len(all.Value()) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(all.Value()))
		}
	})
}

func TestSQLiteTaskRepository_Delete(t *testing.T) {
	t.Run("deletes a task", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewTaskRepository(db)

		r := task.NewTask("To delete")
		repo.Save(r.Value())

		delResult := repo.Delete(r.Value().ID())
		if delResult.IsErr() {
			t.Fatalf("delete failed: %s", delResult.Error())
		}

		found := repo.FindByID(r.Value().ID())
		if found.IsOk() {
			t.Error("task should have been deleted")
		}
	})
}

func TestSQLiteTaskRepository_SaveUpdatesExisting(t *testing.T) {
	t.Run("save overwrites existing task", func(t *testing.T) {
		db := setupTestDB(t)
		repo := sqliterepo.NewTaskRepository(db)

		r := task.NewTask("Original")
		tk := r.Value()
		repo.Save(tk)

		tk.UpdateTitle("Updated")
		repo.Save(tk)

		found := repo.FindByID(tk.ID())
		if found.Value().Title() != "Updated" {
			t.Errorf("expected title 'Updated', got '%s'", found.Value().Title())
		}
	})
}
