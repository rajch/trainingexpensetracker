package memrepository

import (
	"errors"
	"sync"
	"time"

	"github.com/rajch/trainingexpensetracker/internal/model"
)

type MemExpenseRepository struct {
	mu       sync.RWMutex
	expenses []model.Expense
	nextId   int
}

func NewMemExpenseRepository() *MemExpenseRepository {
	return &MemExpenseRepository{
		expenses: make([]model.Expense, 0),
		nextId:   1,
	}
}

func (r *MemExpenseRepository) New(date time.Time, description string, amount float64) (model.Expense, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	expense := model.Expense{
		Id:          r.nextId,
		Date:        date,
		Description: description,
		Amount:      amount,
	}
	r.expenses = append(r.expenses, expense)
	r.nextId++

	return expense, nil
}

func (r *MemExpenseRepository) Get(id int) (model.Expense, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, e := range r.expenses {
		if e.Id == id {
			return e, nil
		}
	}
	return model.Expense{}, errors.New("expense not found")
}

func (r *MemExpenseRepository) GetAll() ([]model.Expense, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to avoid external modification issues if the slice is modified outside
	expenses := make([]model.Expense, len(r.expenses))
	copy(expenses, r.expenses)
	return expenses, nil
}

func (r *MemExpenseRepository) Update(expense model.Expense) (model.Expense, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, e := range r.expenses {
		if e.Id == expense.Id {
			r.expenses[i] = expense
			return expense, nil
		}
	}
	return model.Expense{}, errors.New("expense not found")
}

func (r *MemExpenseRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, e := range r.expenses {
		if e.Id == id {
			r.expenses = append(r.expenses[:i], r.expenses[i+1:]...)
			return nil
		}
	}
	return errors.New("expense not found")
}
