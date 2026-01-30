package sqliterepository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rajch/trainingexpensetracker/internal/model"
	_ "modernc.org/sqlite"
)

type SqliteExpenseRepository struct {
	db *sql.DB
}

func NewSqliteExpenseRepository(dbPath string) (*SqliteExpenseRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			description TEXT,
			amount REAL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SqliteExpenseRepository{db: db}, nil
}

func (r *SqliteExpenseRepository) New(date time.Time, description string, amount float64) (model.Expense, error) {
	res, err := r.db.Exec("INSERT INTO expenses (date, description, amount) VALUES (?, ?, ?)",
		date.Format(time.RFC3339), description, amount)
	if err != nil {
		return model.Expense{}, fmt.Errorf("failed to insert expense: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Expense{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return model.Expense{
		Id:          int(id),
		Date:        date,
		Description: description,
		Amount:      amount,
	}, nil
}

func (r *SqliteExpenseRepository) Get(id int) (model.Expense, error) {
	var e model.Expense
	var dateStr string
	err := r.db.QueryRow("SELECT id, date, description, amount FROM expenses WHERE id = ?", id).
		Scan(&e.Id, &dateStr, &e.Description, &e.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Expense{}, fmt.Errorf("expense not found")
		}
		return model.Expense{}, fmt.Errorf("failed to get expense: %w", err)
	}

	e.Date, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return model.Expense{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return e, nil
}

func (r *SqliteExpenseRepository) GetAll() ([]model.Expense, error) {
	rows, err := r.db.Query("SELECT id, date, description, amount FROM expenses")
	if err != nil {
		return nil, fmt.Errorf("failed to query expenses: %w", err)
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var e model.Expense
		var dateStr string
		if err := rows.Scan(&e.Id, &dateStr, &e.Description, &e.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}
		e.Date, _ = time.Parse(time.RFC3339, dateStr)
		expenses = append(expenses, e)
	}

	return expenses, nil
}

func (r *SqliteExpenseRepository) Update(expense model.Expense) (model.Expense, error) {
	_, err := r.db.Exec("UPDATE expenses SET date = ?, description = ?, amount = ? WHERE id = ?",
		expense.Date.Format(time.RFC3339), expense.Description, expense.Amount, expense.Id)
	if err != nil {
		return model.Expense{}, fmt.Errorf("failed to update expense: %w", err)
	}

	// Double check if rows were affected could be added here, but simplicity first.
	return expense, nil
}

func (r *SqliteExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}
	return nil
}

func (r *SqliteExpenseRepository) Close() error {
	return r.db.Close()
}
