package app

import (
	"errors"
	"strings"
	"time"

	"github.com/rajch/trainingexpensetracker/internal/model"
)

type App struct {
	repo model.ExpenseRepository
}

func New(repo model.ExpenseRepository) *App {
	return &App{
		repo: repo,
	}
}

func (a *App) validate(description string, amount float64) error {
	if strings.TrimSpace(description) == "" {
		return errors.New("description cannot be empty")
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

func (a *App) AddExpense(date time.Time, description string, amount float64) (model.Expense, error) {
	if err := a.validate(description, amount); err != nil {
		return model.Expense{}, err
	}
	return a.repo.New(date, description, amount)
}

func (a *App) UpdateExpense(expense model.Expense) (model.Expense, error) {
	if err := a.validate(expense.Description, expense.Amount); err != nil {
		return model.Expense{}, err
	}
	return a.repo.Update(expense)
}

func (a *App) GetAllExpenses() ([]model.Expense, error) {
	return a.repo.GetAll()
}

func (a *App) GetExpenseByID(id int) (model.Expense, error) {
	return a.repo.Get(id)
}

func (a *App) DeleteExpense(id int) error {
	return a.repo.Delete(id)
}
