package model

import "time"

type ExpenseRepository interface {
	New(date time.Time, description string, amount float64) (Expense, error)
	Get(id int) (Expense, error)
	GetAll() ([]Expense, error)
	Update(expense Expense) (Expense, error)
	Delete(id int) error
}
