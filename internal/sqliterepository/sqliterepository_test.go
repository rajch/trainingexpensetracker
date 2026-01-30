package sqliterepository

import (
	"os"
	"testing"

	"github.com/rajch/trainingexpensetracker/internal/model"
)

func TestSqliteExpenseRepository(t *testing.T) {
	dbPath := "test_repo.db"

	// Cleanup from previous runs if any
	os.Remove(dbPath)
	defer os.Remove(dbPath)

	repo, err := NewSqliteExpenseRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	model.RunRepositoryTests(t, repo)
}
