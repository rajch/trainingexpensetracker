package app

import (
	"testing"
	"time"

	"github.com/rajch/trainingexpensetracker/internal/memrepository"
)

func TestAppValidations(t *testing.T) {
	repo := memrepository.NewMemExpenseRepository()
	application := New(repo)

	t.Run("AddExpenseValid", func(t *testing.T) {
		_, err := application.AddExpense(time.Now(), "Coffee", 5.0)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("AddExpenseEmptyDescription", func(t *testing.T) {
		_, err := application.AddExpense(time.Now(), "", 5.0)
		if err == nil || err.Error() != "description cannot be empty" {
			t.Errorf("Expected 'description cannot be empty', got %v", err)
		}
	})

	t.Run("AddExpenseZeroAmount", func(t *testing.T) {
		_, err := application.AddExpense(time.Now(), "Coffee", 0.0)
		if err == nil || err.Error() != "amount must be greater than zero" {
			t.Errorf("Expected 'amount must be greater than zero', got %v", err)
		}
	})

	t.Run("AddExpenseNegativeAmount", func(t *testing.T) {
		_, err := application.AddExpense(time.Now(), "Coffee", -1.0)
		if err == nil || err.Error() != "amount must be greater than zero" {
			t.Errorf("Expected 'amount must be greater than zero', got %v", err)
		}
	})

	t.Run("UpdateExpenseValid", func(t *testing.T) {
		e, _ := application.AddExpense(time.Now(), "Lunch", 10.0)
		e.Description = "Dinner"
		e.Amount = 15.0
		_, err := application.UpdateExpense(e)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("UpdateExpenseInvalid", func(t *testing.T) {
		e, _ := application.AddExpense(time.Now(), "Lunch", 10.0)
		e.Description = ""
		_, err := application.UpdateExpense(e)
		if err == nil {
			t.Errorf("Expected error for empty description in update")
		}
	})

	t.Run("DeleteExpense", func(t *testing.T) {
		e, _ := application.AddExpense(time.Now(), "Lunch", 10.0)
		err := application.DeleteExpense(e.Id)
		if err != nil {
			t.Errorf("Expected no error on delete, got %v", err)
		}
	})
}
