package model

import (
	"testing"
	"time"
)

// RunRepositoryTests is a helper function that runs a suite of tests against any ExpenseRepository implementation.
func RunRepositoryTests(t *testing.T, repo ExpenseRepository) {
	t.Run("CreateAndGet", func(t *testing.T) {
		date := time.Now().Truncate(time.Second) // SQLite might lose some precision if not careful, RFC3339 helps
		description := "Test Expense"
		amount := 10.50

		e, err := repo.New(date, description, amount)
		if err != nil {
			t.Fatalf("Failed to create expense: %v", err)
		}

		if e.Id == 0 {
			t.Errorf("Expected non-zero ID")
		}

		got, err := repo.Get(e.Id)
		if err != nil {
			t.Fatalf("Failed to get expense: %v", err)
		}

		if got.Description != description {
			t.Errorf("Expected description %s, got %s", description, got.Description)
		}
		if got.Amount != amount {
			t.Errorf("Expected amount %v, got %v", amount, got.Amount)
		}
		// Compare times carefully due to serialization
		if !got.Date.Equal(date) && got.Date.Format(time.RFC3339) != date.Format(time.RFC3339) {
			t.Errorf("Expected date %v, got %v", date, got.Date)
		}
	})

	t.Run("Update", func(t *testing.T) {
		e, _ := repo.New(time.Now(), "Update Me", 1.0)
		e.Description = "Updated"
		e.Amount = 2.0

		updated, err := repo.Update(e)
		if err != nil {
			t.Fatalf("Failed to update expense: %v", err)
		}

		if updated.Description != "Updated" || updated.Amount != 2.0 {
			t.Errorf("Update returned wrong values: %+v", updated)
		}

		got, _ := repo.Get(e.Id)
		if got.Description != "Updated" {
			t.Errorf("Database has wrong value after update: %s", got.Description)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		// Clear or at least ensure we have some
		repo.New(time.Now(), "Exp 1", 1.0)
		repo.New(time.Now(), "Exp 2", 2.0)

		all, err := repo.GetAll()
		if err != nil {
			t.Fatalf("Failed to get all expenses: %v", err)
		}

		if len(all) < 2 {
			t.Errorf("Expected at least 2 expenses, got %d", len(all))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		e, _ := repo.New(time.Now(), "Delete Me", 1.0)

		err := repo.Delete(e.Id)
		if err != nil {
			t.Fatalf("Failed to delete expense: %v", err)
		}

		_, err = repo.Get(e.Id)
		if err == nil {
			t.Errorf("Expected error when getting deleted expense, got nil")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := repo.Get(999999)
		if err == nil {
			t.Errorf("Expected error for non-existent ID, got nil")
		}
	})
}
